package consume

import (
	"errors"
	"testing"
	"time"

	"github.com/allegedlyreliable/sqlstreams/pkg/common"
	"github.com/allegedlyreliable/sqlstreams/pkg/consume/messageconsumer/controller/datastore"
	"github.com/allegedlyreliable/sqlstreams/pkg/stream"
)

// invariant (per-key order): a successor waits for its predecessor's retry
// to succeed, then starts at attempt 0. Deferring that successor again does
// not spend an attempt before its handler runs.
func TestOrderedSuccessorStartsAtAttemptZeroAfterPredecessorRetry(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceKeyedMessages(t, groups, consumer, "userkey-123", 2)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			outcomes := []datastore.Outcome{
				{MessageId: 1, MessageKey: "userkey-123", Concurrency: common.ConcurrencyOrdered, Kind: datastore.OutcomeException, Err: "first attempt failed"},
				{MessageId: 2, MessageKey: "userkey-123", Concurrency: common.ConcurrencyOrdered, Kind: datastore.OutcomeDeferred},
			}
			if err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, outcomes, 0, mode); err != nil {
				t.Fatal(err)
			}

			// test
			predecessor, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(predecessor) != 1 || predecessor[0].MessageId != 1 || predecessor[0].Attempts != 1 {
				t.Fatalf("Claim after first failure = %+v, %v; want message 1 at attempt 1", predecessor, err)
			}

			// test
			if err := exceptions.RecordSuccess(t.Context(), &predecessor[0], mode, nil); err != nil {
				t.Fatal(err)
			}
			successor, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(successor) != 1 || successor[0].MessageId != 2 || successor[0].Attempts != 0 {
				t.Fatalf("Claim after predecessor success = %+v, %v; want message 2 at attempt 0", successor, err)
			}

			// test
			if err := exceptions.RecordDeferred(t.Context(), &successor[0], common.ConcurrencyOrdered, mode); err != nil {
				t.Fatal(err)
			}
			successor, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(successor) != 1 || successor[0].Attempts != 0 {
				t.Fatalf("Claim after another key deferral = %+v, %v; want attempt 0", successor, err)
			}
		})
	}
}

// invariant (retry budget): ready and deferred rows at the final retry can
// still be claimed at that attempt. An expired inflight row at the same
// limit cannot start another attempt and is marked dead instead.
func TestFinalRetryCanBeClaimedButCannotBeReclaimedAfterExpiry(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 3)
			insertException(t, groups, consumer, 1, "ready", -time.Hour, 0)
			insertException(t, groups, consumer, 2, "deferred", 0, 0)
			insertException(t, groups, consumer, 3, "inflight", 0, -time.Hour)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET attempts = 3 WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)
			killed, killErr := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 3, mode)

			// verify
			if err != nil || len(claimed) != 2 || claimed[0].MessageId != 1 || claimed[1].MessageId != 2 || claimed[0].Attempts != 3 || claimed[1].Attempts != 3 {
				t.Fatalf("Claim at retry limit = %+v, %v; want ready and deferred messages at attempt 3", claimed, err)
			}
			if killErr != nil || killed != 1 {
				t.Fatalf("Kill at retry limit = %d, %v; want only the expired message killed", killed, killErr)
			}
		})
	}
}

// behavior: delay, failure, and lease expiry each advance the next attempt,
// but delays spend no failure retries. The first failure uses initial
// backoff; later failures use the retry curve, and history names the attempt
// that ended rather than the next one.
func TestExceptionAttemptsAdvanceOnOutcomesAndExpiryWithoutChargingDelays(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			insertException(t, groups, consumer, 1, "deferred", 0, 0)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			retry := (&common.RetryPolicy{MaxRetries: 3, BaseDelay: time.Second}).WithDefaults()
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)
			if err != nil || len(claimed) != 1 {
				t.Fatalf("Claim during setup = %+v, %v; want one message", claimed, err)
			}

			// test
			if err := exceptions.RecordDelayed(t.Context(), 0, &claimed[0], errors.New("try later"), mode, nil); err != nil {
				t.Fatal(err)
			}
			claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1 || claimed[0].Delays != 1 {
				t.Fatalf("Claim after Delay = %+v, %v; want attempt 1, delays 1", claimed, err)
			}

			// test: the first failure follows a delay, with no prior failures
			if err := exceptions.RecordFailure(t.Context(), retry, time.Hour, &claimed[0], errors.New("first failure"), mode, nil); err != nil {
				t.Fatal(err)
			}

			// verify
			var attempts, delays int
			var backoff float64
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT attempts, delays, extract(epoch FROM can_run_after - updated_at) FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&attempts, &delays, &backoff); err != nil {
				t.Fatal(err)
			}
			if attempts != 2 || delays != 1 || backoff != time.Hour.Seconds() {
				t.Fatalf("first failure: next attempt, delays, backoff = %d, %d, %g; want 2, 1, 3600", attempts, delays, backoff)
			}

			// test: the next failure uses the retry curve's first backoff
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET can_run_after = now() WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)
			if err != nil || len(claimed) != 1 {
				t.Fatalf("Claim after first failure = %+v, %v; want one message", claimed, err)
			}
			if err := exceptions.RecordFailure(t.Context(), retry, time.Hour, &claimed[0], errors.New("second failure"), mode, nil); err != nil {
				t.Fatal(err)
			}

			// verify
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT attempts, delays, extract(epoch FROM can_run_after - updated_at) FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&attempts, &delays, &backoff); err != nil {
				t.Fatal(err)
			}
			if attempts != 3 || delays != 1 || backoff != time.Second.Seconds() {
				t.Fatalf("second failure: next attempt, delays, backoff = %d, %d, %g; want 3, 1, 1", attempts, delays, backoff)
			}

			// test: abandon attempt 3, then claim the last allowed retry
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET can_run_after = now() WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)
			if err != nil || len(claimed) != 1 {
				t.Fatalf("Claim after second failure = %+v, %v; want one message", claimed, err)
			}
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET lease_expires_at = now() - interval '1 second' WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 1 || claimed[0].Attempts != 4 || claimed[0].Delays != 1 {
				t.Fatalf("Claim after expiry = %+v, %v; want attempt 4, delays 1 (retry 3)", claimed, err)
			}

			// test: expiry of the final retry cannot produce another claim
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET lease_expires_at = now() - interval '1 second' WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)
			killed, killErr := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 3, mode)

			// verify
			if err != nil || len(claimed) != 0 || killErr != nil || killed != 1 {
				t.Fatalf("Claim, Kill after final expiry = %+v, %v, %d, %v; want no claim and one dead letter", claimed, err, killed, killErr)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "delayed:0,failure:1,failure:2,expired:3,killed:4"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("delivery history = %q, want %q", history, want)
			}
		})
	}
}

// behavior: every cursor outcome describes attempt 0. Failure and delay
// store next attempt 1; deferral and terminal outcomes retain 0. Success
// and supersession leave no exception row, and enabled history records 0.
func TestCursorOutcomesInitializeAttemptsAndHistory(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			for _, outcome := range []struct {
				name     string
				kind     datastore.OutcomeKind
				status   string
				attempts int
				delays   int
				logged   string
			}{
				{"success", datastore.OutcomeSuccess, "", 0, 0, "success:0"},
				{"failure", datastore.OutcomeException, "ready", 1, 0, "failure:0"},
				{"terminal", datastore.OutcomeTerminal, "dead", 0, 0, "failure:0"},
				{"deferred", datastore.OutcomeDeferred, "deferred", 0, 0, "deferred:0"},
				{"delayed", datastore.OutcomeDelayed, "ready", 1, 1, "delayed:0"},
				{"superseded", datastore.OutcomeSuperseded, "", 0, 0, "superseded:0"},
			} {
				t.Run(outcome.name, func(t *testing.T) {
					// setup
					groups, consumer := newMessageConsumerDatastore(t)
					produceMessages(t, groups, consumer, 1)
					claimed := claimRange(t, groups, consumer)
					exceptions := newExceptionConsumerDatastore(t, groups)
					queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
					logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
					outcomes := []datastore.Outcome{{MessageId: 1, Kind: outcome.kind, Err: "handler outcome"}}
					// Cursor callers collect success outcomes only in all mode;
					// otherwise success is represented by closing the range alone.
					if outcome.kind == datastore.OutcomeSuccess && mode != stream.DeliveryLogModeAll {
						outcomes = nil
					}

					// test
					err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, claimed.Lease.Token, outcomes, 0, mode)

					// verify
					if err != nil {
						t.Fatal(err)
					}
					var count, attempts, delays int
					var status, history string
					if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT count(*), coalesce(max(status), ''), coalesce(max(attempts), 0), coalesce(max(delays), 0) FROM "+queue+" WHERE consumer_group_id = $1", consumer.Id).Scan(&count, &status, &attempts, &delays); err != nil {
						t.Fatal(err)
					}
					wantCount := 1
					if outcome.status == "" {
						wantCount = 0
					}
					if count != wantCount || status != outcome.status || attempts != outcome.attempts || delays != outcome.delays {
						t.Fatalf("Commit(%s): rows, status, attempts, delays = %d, %q, %d, %d; want %d, %q, %d, %d", outcome.name, count, status, attempts, delays, wantCount, outcome.status, outcome.attempts, outcome.delays)
					}
					if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
						t.Fatal(err)
					}
					wantHistory := outcome.logged
					if mode == stream.DeliveryLogModeOff || mode == stream.DeliveryLogModeFailures && outcome.kind == datastore.OutcomeSuccess {
						wantHistory = ""
					}
					if history != wantHistory {
						t.Errorf("Commit(%s) history = %q, want %q", outcome.name, history, wantHistory)
					}

					// test: the next claim sees the number already stored by the outcome
					next, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

					// verify
					if err != nil {
						t.Fatal(err)
					}
					if outcome.status == "ready" || outcome.status == "deferred" {
						if len(next) != 1 || next[0].Attempts != outcome.attempts || next[0].Delays != outcome.delays {
							t.Fatalf("Claim after %s = %+v, want attempt %d, delays %d", outcome.name, next, outcome.attempts, outcome.delays)
						}
					} else if len(next) != 0 {
						t.Fatalf("Claim after %s = %+v, want no delivery", outcome.name, next)
					}
				})
			}
		})
	}
}

// behavior: an initial failure followed by two failed retries delivers
// attempts 1, 2, and 3 in order. Success on the final retry removes the
// exception, and history records each outcome at its delivery's attempt.
func TestConsecutiveFailuresReachTheFinalRetryWithoutSkippingAnAttempt(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			retry := (&common.RetryPolicy{MaxRetries: 3}).WithDefaults()

			// test
			err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, []datastore.Outcome{{MessageId: 1, Kind: datastore.OutcomeException, Err: "initial failure"}}, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			for attempt := 1; attempt <= 3; attempt++ {
				// test
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != attempt || claimed[0].Delays != 0 {
					t.Fatalf("Claim for retry %d = %+v, %v; want attempt %d, delays 0", attempt, claimed, err, attempt)
				}

				// test
				if attempt == 3 {
					err = exceptions.RecordSuccess(t.Context(), &claimed[0], mode, nil)
				} else {
					err = exceptions.RecordFailure(t.Context(), retry, time.Hour, &claimed[0], errors.New("retry failed"), mode, nil)
				}

				// verify
				if err != nil {
					t.Fatal(err)
				}
				var count, nextAttempt int
				if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT count(*), coalesce(max(attempts), 0) FROM "+queue+" WHERE consumer_group_id = $1", consumer.Id).Scan(&count, &nextAttempt); err != nil {
					t.Fatal(err)
				}
				if attempt < 3 && (count != 1 || nextAttempt != attempt+1) {
					t.Fatalf("RecordFailure at %d: rows, next attempt = %d, %d; want 1, %d", attempt, count, nextAttempt, attempt+1)
				}
				if attempt == 3 && count != 0 {
					t.Fatalf("RecordSuccess at final retry: rows = %d, want 0", count)
				}
				if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET can_run_after = now() WHERE consumer_group_id = $1", consumer.Id); err != nil {
					t.Fatal(err)
				}
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "failure:0,failure:1,failure:2"
			if mode == stream.DeliveryLogModeAll {
				want += ",success:3"
			} else if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("retry history = %q, want %q", history, want)
			}
		})
	}
}

// invariant (retry budget): requested delays remain claimable on the final
// failure retry because attempts and delays advance together. A terminal
// outcome keeps the current attempt in the dead row and its history.
func TestDelaysOnTheFinalRetryPreserveTheFailureBudget(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)

			// test: a one-retry budget has been spent when the first retry starts
			err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, []datastore.Outcome{{MessageId: 1, Kind: datastore.OutcomeException, Err: "initial failure"}}, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			for delays := 0; delays <= 3; delays++ {
				// test
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 1, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1+delays || claimed[0].Delays != delays {
					t.Fatalf("Claim after %d delays on final retry = %+v, %v; want attempt %d, delays %d", delays, claimed, err, 1+delays, delays)
				}

				// test: delays request another delivery; a terminal outcome does not
				if delays < 3 {
					err = exceptions.RecordDelayed(t.Context(), 0, &claimed[0], errors.New("not ready"), mode, nil)
				} else {
					err = exceptions.RecordTerminal(t.Context(), &claimed[0], errors.New("final retry failed"), mode, nil)
				}

				// verify
				if err != nil {
					t.Fatal(err)
				}
			}
			var status, history string
			var attempts, delays int
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT status, attempts, delays FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&status, &attempts, &delays); err != nil {
				t.Fatal(err)
			}
			if status != "dead" || attempts != 4 || delays != 3 {
				t.Fatalf("terminal final retry: status, attempts, delays = %q, %d, %d; want dead, 4, 3", status, attempts, delays)
			}
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "failure:0,delayed:1,delayed:2,delayed:3,failure:4"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("final retry history = %q, want %q", history, want)
			}

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 1, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 0 {
				t.Fatalf("Claim after terminal final retry = %+v, %v; want no delivery", claimed, err)
			}
		})
	}
}

// behavior: repeated key deferrals preserve an already nonzero attempt and
// its delay count. Supersession retains those counters, prevents another
// claim, and logs the same attempt as the preceding deferrals.
func TestRepeatedDeferralsAndSupersessionPreserveANonzeroAttempt(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceKeyedMessages(t, groups, consumer, "invoice-42", 1)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)

			// test
			err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, []datastore.Outcome{{MessageId: 1, MessageKey: "invoice-42", Kind: datastore.OutcomeDelayed}}, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			for deferrals := 0; deferrals <= 3; deferrals++ {
				// test
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1 || claimed[0].Delays != 1 {
					t.Fatalf("Claim after %d deferrals = %+v, %v; want attempt 1, delays 1", deferrals, claimed, err)
				}

				// test
				if deferrals < 3 {
					err = exceptions.RecordDeferred(t.Context(), &claimed[0], common.ConcurrencyExclusive, mode)
				} else {
					err = exceptions.RecordSuperseded(t.Context(), &claimed[0], mode)
				}

				// verify
				if err != nil {
					t.Fatal(err)
				}
			}
			var attempts, delays int
			var status, history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT status, attempts, delays FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&status, &attempts, &delays); err != nil {
				t.Fatal(err)
			}
			if status != "superseded" || attempts != 1 || delays != 1 {
				t.Fatalf("superseded delivery: status, attempts, delays = %q, %d, %d; want superseded, 1, 1", status, attempts, delays)
			}
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "delayed:0,deferred:1,deferred:1,deferred:1,superseded:1"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("deferral history = %q, want %q", history, want)
			}

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 0 {
				t.Fatalf("Claim after supersession = %+v, %v; want no delivery", claimed, err)
			}
		})
	}
}

// invariant (one live lease): after expiry and takeover, the old holder
// cannot record outcomes or renew the lease. The current holder's failure
// advances attempts once; repeated or late outcomes change neither counters
// nor history.
func TestStaleAndRepeatedOutcomesCannotAdvanceAttemptsTwice(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			insertException(t, groups, consumer, 1, "deferred", 0, 0)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			retry := (&common.RetryPolicy{MaxRetries: 3}).WithDefaults()

			// test
			original, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(original) != 1 || original[0].Attempts != 0 {
				t.Fatalf("initial Claim = %+v, %v; want attempt 0", original, err)
			}

			// test: another consumer takes over the expired delivery
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET lease_expires_at = now() - interval '1 second' WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			replacement, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(replacement) != 1 || replacement[0].Attempts != 1 {
				t.Fatalf("Claim after expiry = %+v, %v; want attempt 1", replacement, err)
			}

			// test: late results from the old holder must have no effect
			failureErr := exceptions.RecordFailure(t.Context(), retry, time.Hour, &original[0], errors.New("late failure"), mode, nil)
			delayErr := exceptions.RecordDelayed(t.Context(), 0, &original[0], errors.New("late delay"), mode, nil)
			deferredErr := exceptions.RecordDeferred(t.Context(), &original[0], common.ConcurrencyExclusive, mode)
			supersededErr := exceptions.RecordSuperseded(t.Context(), &original[0], mode)
			terminalErr := exceptions.RecordTerminal(t.Context(), &original[0], errors.New("late terminal"), mode, nil)
			successErr := exceptions.RecordSuccess(t.Context(), &original[0], mode, nil)
			renewed, renewErr := exceptions.RenewLease(t.Context(), &original[0], time.Hour)

			// verify
			if !errors.Is(failureErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordFailure = %v, want ErrLeaseLost", failureErr)
			}
			if !errors.Is(delayErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordDelayed = %v, want ErrLeaseLost", delayErr)
			}
			if !errors.Is(deferredErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordDeferred = %v, want ErrLeaseLost", deferredErr)
			}
			if !errors.Is(supersededErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordSuperseded = %v, want ErrLeaseLost", supersededErr)
			}
			if !errors.Is(terminalErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordTerminal = %v, want ErrLeaseLost", terminalErr)
			}
			if !errors.Is(successErr, common.ErrLeaseLost) {
				t.Errorf("stale RecordSuccess = %v, want ErrLeaseLost", successErr)
			}
			if renewErr != nil || renewed {
				t.Errorf("stale RenewLease = %t, %v; want false, nil", renewed, renewErr)
			}
			var attempts, delays int
			var status string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT status, attempts, delays FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&status, &attempts, &delays); err != nil {
				t.Fatal(err)
			}
			if status != "inflight" || attempts != 1 || delays != 0 {
				t.Fatalf("after stale outcomes: status, attempts, delays = %q, %d, %d; want inflight, 1, 0", status, attempts, delays)
			}

			// test: recording the current holder's failure twice advances only once
			failureErr = exceptions.RecordFailure(t.Context(), retry, time.Hour, &replacement[0], errors.New("current failure"), mode, nil)
			repeatedErr := exceptions.RecordFailure(t.Context(), retry, time.Hour, &replacement[0], errors.New("repeated failure"), mode, nil)
			delayErr = exceptions.RecordDelayed(t.Context(), 0, &replacement[0], errors.New("late delay after failure"), mode, nil)

			// verify
			if failureErr != nil || !errors.Is(repeatedErr, common.ErrLeaseLost) || !errors.Is(delayErr, common.ErrLeaseLost) {
				t.Fatalf("current, repeated failure, late delay = %v, %v, %v; want nil, ErrLeaseLost, ErrLeaseLost", failureErr, repeatedErr, delayErr)
			}
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT status, attempts, delays FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&status, &attempts, &delays); err != nil {
				t.Fatal(err)
			}
			if status != "ready" || attempts != 2 || delays != 0 {
				t.Fatalf("after repeated failure: status, attempts, delays = %q, %d, %d; want ready, 2, 0", status, attempts, delays)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "expired:0,failure:1"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("history after stale and repeated outcomes = %q, want %q", history, want)
			}
		})
	}
}

// invariant (retry budget): renewal, polling, and the kill backstop leave a
// live attempt unchanged. Each expired lease permits exactly the next
// attempt until the budget is exhausted; killing the final attempt twice
// writes only one killed outcome.
func TestRenewalsAndPollingPreserveAttemptsUntilEachLeaseExpires(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			insertException(t, groups, consumer, 1, "deferred", 0, 0)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 2, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 1 || claimed[0].Attempts != 0 {
				t.Fatalf("initial Claim = %+v, %v; want attempt 0", claimed, err)
			}
			for attempt := 0; attempt <= 2; attempt++ {
				// test
				renewed, renewErr := exceptions.RenewLease(t.Context(), &claimed[0], time.Hour)
				polled, pollErr := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 2, time.Minute, mode)
				killed, killErr := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 2, mode)

				// verify
				if renewErr != nil || !renewed || pollErr != nil || len(polled) != 0 || killErr != nil || killed != 0 {
					t.Fatalf("live attempt %d: RenewLease, Claim, Kill = %t, %v, %+v, %v, %d, %v; want renewed, no claim, no kill", attempt, renewed, renewErr, polled, pollErr, killed, killErr)
				}
				var stored int
				if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT attempts FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&stored); err != nil {
					t.Fatal(err)
				}
				if stored != attempt {
					t.Fatalf("attempt after renewal and polling = %d, want %d", stored, attempt)
				}

				// test
				if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET lease_expires_at = now() - interval '1 second' WHERE consumer_group_id = $1", consumer.Id); err != nil {
					t.Fatal(err)
				}
				claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 2, time.Minute, mode)

				// verify
				if err != nil {
					t.Fatal(err)
				}
				if attempt < 2 && (len(claimed) != 1 || claimed[0].Attempts != attempt+1 || claimed[0].Delays != 0) {
					t.Fatalf("Claim after expiry of %d = %+v, want attempt %d, delays 0", attempt, claimed, attempt+1)
				}
				if attempt == 2 && len(claimed) != 0 {
					t.Fatalf("Claim after final retry expired = %+v, want no claim", claimed)
				}
			}

			// test
			killed, err := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 2, mode)
			killedAgain, repeatedErr := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 2, mode)

			// verify
			if err != nil || killed != 1 || repeatedErr != nil || killedAgain != 0 {
				t.Fatalf("Kill, repeated Kill = %d, %v, %d, %v; want 1, nil, 0, nil", killed, err, killedAgain, repeatedErr)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "expired:0,expired:1,killed:2"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("expiry history = %q, want %q", history, want)
			}
		})
	}
}

// behavior: a deferred delivery with zero retries still gets attempt 0 and
// requested later runs. Each recorded delay advances attempts and delays
// once, but an expired lease needs a failure retry and is marked dead.
func TestDeferredDeliveryWithNoRetriesCanDelayButCannotRecoverAnExpiredLease(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)

			// test
			err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, []datastore.Outcome{{MessageId: 1, Kind: datastore.OutcomeDeferred}}, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			for delays := 0; delays <= 2; delays++ {
				// test: MaxRetries zero still permits the first run and requested delays
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 0, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != delays || claimed[0].Delays != delays {
					t.Fatalf("Claim with no retries after %d delays = %+v, %v; want attempt and delays %d", delays, claimed, err, delays)
				}
				if delays == 2 {
					break
				}

				// test
				err = exceptions.RecordDelayed(t.Context(), 0, &claimed[0], errors.New("not ready"), mode, nil)
				repeatedErr := exceptions.RecordDelayed(t.Context(), 0, &claimed[0], errors.New("repeated delay"), mode, nil)

				// verify
				if err != nil || !errors.Is(repeatedErr, common.ErrLeaseLost) {
					t.Fatalf("RecordDelayed, repeated RecordDelayed = %v, %v; want nil, ErrLeaseLost", err, repeatedErr)
				}
			}

			// test
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET lease_expires_at = now() - interval '1 second' WHERE consumer_group_id = $1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			claimed, claimErr := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 0, time.Minute, mode)
			killed, killErr := exceptions.Kill(t.Context(), consumer.StreamId, consumer.Id, 0, mode)

			// verify
			if claimErr != nil || len(claimed) != 0 || killErr != nil || killed != 1 {
				t.Fatalf("Claim, Kill after expiry with no retries = %+v, %v, %d, %v; want no claim and one dead letter", claimed, claimErr, killed, killErr)
			}
			var status, history string
			var attempts, delays int
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT status, attempts, delays FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&status, &attempts, &delays); err != nil {
				t.Fatal(err)
			}
			if status != "dead" || attempts != 2 || delays != 2 {
				t.Fatalf("expired delivery with no retries: status, attempts, delays = %q, %d, %d; want dead, 2, 2", status, attempts, delays)
			}
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "deferred:0,delayed:0,delayed:1,killed:2"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("history with no retries = %q, want %q", history, want)
			}
		})
	}
}

// invariant (crash consistency): a partial commit preserves each recorded
// outcome's next attempt when the remaining range is quarantined. Only the
// unrecorded suffix starts a fresh per-message budget at 0, whose first
// failure still uses initial backoff; the old range holder cannot overwrite it.
func TestPartialCommitAndQuarantineKeepSeparateMessageAttemptBudgets(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 4)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			outcomes := []datastore.Outcome{
				{MessageId: 1, Kind: datastore.OutcomeException, Err: "first failure"},
				{MessageId: 2, Kind: datastore.OutcomeDelayed},
				{MessageId: 3, Kind: datastore.OutcomeDeferred},
			}

			// test
			err := groups.PartialCommit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, 3, outcomes, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}

			// test: a crash quarantines only the unrecorded suffix, message 4
			expireClaimLease(t, groups, consumer)
			quarantined, err := groups.ClaimMessagesWithCursor(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 1, time.Minute, mode)

			// verify
			if err != nil || quarantined == nil || !quarantined.Quarantined || quarantined.Lease.Low != 3 || quarantined.Lease.High != 4 {
				t.Fatalf("ClaimMessagesWithCursor after partial commit and crash = %+v, %v; want quarantined range (3, 4]", quarantined, err)
			}

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)
			staleErr := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, outcomes, 0, mode)

			// verify
			if err != nil || len(claimed) != 4 {
				t.Fatalf("Claim after partial commit and quarantine = %+v, %v; want four independent deliveries", claimed, err)
			}
			if !errors.Is(staleErr, common.ErrLeaseLost) {
				t.Errorf("Commit from crashed range holder = %v, want ErrLeaseLost", staleErr)
			}
			for i, want := range []struct{ attempts, delays int }{{1, 0}, {1, 1}, {0, 0}, {0, 0}} {
				if claimed[i].MessageId != int64(i+1) || claimed[i].Attempts != want.attempts || claimed[i].Delays != want.delays {
					t.Errorf("Claim message %d = %+v, want attempts %d, delays %d", i+1, claimed[i], want.attempts, want.delays)
				}
			}

			// test: quarantine has not spent message 4's first-failure backoff
			retry := (&common.RetryPolicy{MaxRetries: 3, BaseDelay: time.Second}).WithDefaults()
			err = exceptions.RecordFailure(t.Context(), retry, time.Hour, &claimed[3], errors.New("first independent failure"), mode, nil)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			var attempts int
			var backoff float64
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT attempts, extract(epoch FROM can_run_after - updated_at) FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 4", consumer.Id).Scan(&attempts, &backoff); err != nil {
				t.Fatal(err)
			}
			if attempts != 1 || backoff != time.Hour.Seconds() {
				t.Fatalf("first quarantined failure: next attempt, backoff = %d, %g; want 1, 3600", attempts, backoff)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(message_id || ':' || status || ':' || attempt, ',' ORDER BY message_id, id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			// Quarantine logs the range failure at zero; the independent delivery
			// also starts at zero, because range reclaims have their own budget.
			want := "1:failure:0,2:delayed:0,3:deferred:0,4:failure:0,4:failure:0"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("partial commit and quarantine history = %q, want %q", history, want)
			}
		})
	}
}

// behavior: replaying a crashed cursor range does not spend per-message
// attempts, and surrendering an unstarted range does not spend another
// range reclaim. Its first recorded handler failure logs attempt 0 and
// schedules exception attempt 1 regardless of earlier range replays.
func TestRangeReplaysAndSurrenderDoNotBecomeMessageAttempts(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			for reclaims := 1; reclaims <= 2; reclaims++ {
				// test
				expireClaimLease(t, groups, consumer)
				replayed, err := groups.ClaimMessagesWithCursor(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

				// verify
				if err != nil || replayed == nil || replayed.Quarantined || len(replayed.Messages) != 1 || replayed.Messages[0].Id != 1 || replayed.Lease.Reclaims != reclaims {
					t.Fatalf("range replay %d = %+v, %v; want message 1 with reclaim count %d", reclaims, replayed, err, reclaims)
				}
				rangeClaim = replayed
			}

			// test: surrendering an unstarted range does not use the last reclaim
			err := groups.ForceReclaimRange(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token)

			// verify
			if err != nil {
				t.Fatal(err)
			}

			// test
			surrendered, err := groups.ClaimMessagesWithCursor(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || surrendered == nil || surrendered.Quarantined || surrendered.Lease.Reclaims != 2 || len(surrendered.Messages) != 1 {
				t.Fatalf("claim after surrender = %+v, %v; want the range still at reclaim count 2", surrendered, err)
			}

			// test: the first recorded handler failure still schedules attempt 1
			err = groups.Commit(t.Context(), consumer.StreamId, consumer.Id, surrendered.Lease.Token, []datastore.Outcome{{MessageId: 1, Kind: datastore.OutcomeException, Err: "first recorded failure"}}, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}

			// test
			claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 3, time.Minute, mode)

			// verify
			if err != nil || len(claimed) != 1 || claimed[0].Attempts != 1 || claimed[0].Delays != 0 {
				t.Fatalf("Claim after replayed cursor failure = %+v, %v; want attempt 1, delays 0", claimed, err)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(status || ':' || attempt, ',' ORDER BY id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "failure:0"
			if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("range replay history = %q, want %q", history, want)
			}
		})
	}
}

// invariant (failure backoff): delays advance the delivery number and
// deferrals preserve it, but neither advances the failure backoff curve.
// Interleaved failures still use initial backoff, then exponential backoff
// up to its configured ceiling.
func TestDelaysAndDeferralsDoNotAdvanceTheFailureBackoffCurve(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceMessages(t, groups, consumer, 1)
			insertException(t, groups, consumer, 1, "deferred", 0, 0)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			retry := (&common.RetryPolicy{MaxRetries: 6, BaseDelay: time.Second, Exponent: 2, MaxDelay: 3 * time.Second}).WithDefaults()
			for failures, wantBackoff := range []float64{60, 1, 2, 3, 3} {
				// test: each failure is preceded by one requested delay and a deferral
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 6, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != failures*2 || claimed[0].Delays != failures {
					t.Fatalf("Claim after %d failures and delays = %+v, %v; want attempt %d, delays %d", failures, claimed, err, failures*2, failures)
				}

				// test
				err = exceptions.RecordDelayed(t.Context(), 0, &claimed[0], errors.New("wait for dependency"), mode, nil)

				// verify
				if err != nil {
					t.Fatal(err)
				}

				// test
				claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 6, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != failures*2+1 || claimed[0].Delays != failures+1 {
					t.Fatalf("Claim after delay = %+v, %v; want attempt %d, delays %d", claimed, err, failures*2+1, failures+1)
				}

				// test
				err = exceptions.RecordDeferred(t.Context(), &claimed[0], common.ConcurrencyExclusive, mode)

				// verify
				if err != nil {
					t.Fatal(err)
				}

				// test
				claimed, err = exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 1, 6, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 1 || claimed[0].Attempts != failures*2+1 || claimed[0].Delays != failures+1 {
					t.Fatalf("Claim after deferral = %+v, %v; want unchanged attempt %d, delays %d", claimed, err, failures*2+1, failures+1)
				}

				// test
				err = exceptions.RecordFailure(t.Context(), retry, time.Minute, &claimed[0], errors.New("dependency failed"), mode, nil)

				// verify
				if err != nil {
					t.Fatal(err)
				}
				var attempts, delays int
				var backoff float64
				if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT attempts, delays, extract(epoch FROM can_run_after - updated_at) FROM "+queue+" WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id).Scan(&attempts, &delays, &backoff); err != nil {
					t.Fatal(err)
				}
				if attempts != (failures+1)*2 || delays != failures+1 || backoff != wantBackoff {
					t.Fatalf("failure %d: attempts, delays, backoff = %d, %d, %g; want %d, %d, %g", failures+1, attempts, delays, backoff, (failures+1)*2, failures+1, wantBackoff)
				}
				if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET can_run_after = now() WHERE consumer_group_id = $1", consumer.Id); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}

// invariant (consumer group isolation): each message has its own attempt
// and delay counts in each group. Retrying one message leaves its sibling
// unchanged, and another group's successful deliveries remove only that
// group's exceptions and record its own attempt numbers.
func TestAttemptsBelongToEachMessageAndConsumerGroup(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, emails := newMessageConsumerDatastore(t)
			analytics := registerConsumer(t, groups, emails, "analytics")
			produceMessages(t, groups, emails, 2)
			emailRange := claimRange(t, groups, emails)
			analyticsRange := claimRange(t, groups, analytics)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(emails.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(emails.StreamId)
			retry := (&common.RetryPolicy{MaxRetries: 3}).WithDefaults()

			// test: the same two messages take different paths in each group
			emailErr := groups.Commit(t.Context(), emails.StreamId, emails.Id, emailRange.Lease.Token, []datastore.Outcome{
				{MessageId: 1, Kind: datastore.OutcomeException, Err: "email failed"},
				{MessageId: 2, Kind: datastore.OutcomeDeferred},
			}, 0, mode)
			analyticsErr := groups.Commit(t.Context(), analytics.StreamId, analytics.Id, analyticsRange.Lease.Token, []datastore.Outcome{
				{MessageId: 1, Kind: datastore.OutcomeDeferred},
				{MessageId: 2, Kind: datastore.OutcomeDelayed},
			}, 0, mode)

			// verify
			if emailErr != nil || analyticsErr != nil {
				t.Fatalf("Commit for email, analytics = %v, %v; want nil, nil", emailErr, analyticsErr)
			}

			// test
			emailClaims, emailErr := exceptions.Claim(t.Context(), emails.StreamId, emails.Id, 1, 100, 3, time.Minute, mode)
			analyticsClaims, analyticsErr := exceptions.Claim(t.Context(), analytics.StreamId, analytics.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if emailErr != nil || len(emailClaims) != 2 || emailClaims[0].MessageId != 1 || emailClaims[0].Attempts != 1 || emailClaims[0].Delays != 0 || emailClaims[1].MessageId != 2 || emailClaims[1].Attempts != 0 || emailClaims[1].Delays != 0 {
				t.Fatalf("email Claim = %+v, %v; want message 1 at (1, 0), message 2 at (0, 0)", emailClaims, emailErr)
			}
			if analyticsErr != nil || len(analyticsClaims) != 2 || analyticsClaims[0].MessageId != 1 || analyticsClaims[0].Attempts != 0 || analyticsClaims[0].Delays != 0 || analyticsClaims[1].MessageId != 2 || analyticsClaims[1].Attempts != 1 || analyticsClaims[1].Delays != 1 {
				t.Fatalf("analytics Claim = %+v, %v; want message 1 at (0, 0), message 2 at (1, 1)", analyticsClaims, analyticsErr)
			}

			// test: one group retries while the other resolves both messages
			failureErr := exceptions.RecordFailure(t.Context(), retry, time.Minute, &emailClaims[0], errors.New("email failed again"), mode, nil)
			firstSuccessErr := exceptions.RecordSuccess(t.Context(), &analyticsClaims[0], mode, nil)
			secondSuccessErr := exceptions.RecordSuccess(t.Context(), &analyticsClaims[1], mode, nil)

			// verify
			if failureErr != nil || firstSuccessErr != nil || secondSuccessErr != nil {
				t.Fatalf("email failure, analytics successes = %v, %v, %v; want nil", failureErr, firstSuccessErr, secondSuccessErr)
			}
			var stored string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(message_id || ':' || status || ':' || attempts || ':' || delays, ',' ORDER BY message_id), '') FROM "+queue+" WHERE consumer_group_id = $1", emails.Id).Scan(&stored); err != nil {
				t.Fatal(err)
			}
			if stored != "1:ready:2:0,2:inflight:0:0" {
				t.Fatalf("email queue after analytics success = %q, want 1:ready:2:0,2:inflight:0:0", stored)
			}
			var remaining int
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT count(*) FROM "+queue+" WHERE consumer_group_id = $1", analytics.Id).Scan(&remaining); err != nil {
				t.Fatal(err)
			}
			if remaining != 0 {
				t.Fatalf("analytics queue after success = %d, want 0 rows", remaining)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(message_id || ':' || status || ':' || attempt, ',' ORDER BY message_id, id), '') FROM "+logs+" WHERE consumer_group_id = $1", analytics.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "1:deferred:0,2:delayed:0"
			if mode == stream.DeliveryLogModeAll {
				want = "1:deferred:0,1:success:0,2:delayed:0,2:success:1"
			} else if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("analytics history = %q, want %q", history, want)
			}
		})
	}
}

// invariant (per-key order): polling behind a delayed predecessor spends
// no successor attempts. Terminal and superseded predecessors each release
// the next message at attempt 0, and its history retains that number.
func TestOrderedSuccessorsStayAtZeroThroughDelayTerminalAndSupersession(t *testing.T) {
	for _, mode := range []stream.DeliveryLogMode{stream.DeliveryLogModeOff, stream.DeliveryLogModeFailures, stream.DeliveryLogModeAll} {
		t.Run(string(mode), func(t *testing.T) {
			// setup
			groups, consumer := newMessageConsumerDatastore(t)
			produceKeyedMessages(t, groups, consumer, "order-42", 3)
			rangeClaim := claimRange(t, groups, consumer)
			exceptions := newExceptionConsumerDatastore(t, groups)
			queue := groups.Datastore.Schema + "." + stream.ExceptionQueueTable(consumer.StreamId)
			logs := groups.Datastore.Schema + "." + stream.DeliveryLogTable(consumer.StreamId)
			outcomes := []datastore.Outcome{
				{MessageId: 1, MessageKey: "order-42", Concurrency: common.ConcurrencyOrdered, Kind: datastore.OutcomeDelayed, Delay: time.Hour},
				{MessageId: 2, MessageKey: "order-42", Concurrency: common.ConcurrencyOrdered, Kind: datastore.OutcomeDeferred},
				{MessageId: 3, MessageKey: "order-42", Concurrency: common.ConcurrencyOrdered, Kind: datastore.OutcomeDeferred},
			}

			// test
			err := groups.Commit(t.Context(), consumer.StreamId, consumer.Id, rangeClaim.Lease.Token, outcomes, 0, mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			for range 3 {
				// test: waiting behind a delayed predecessor is not a delivery attempt
				claimed, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

				// verify
				if err != nil || len(claimed) != 0 {
					t.Fatalf("Claim before predecessor is due = %+v, %v; want no deliveries", claimed, err)
				}
			}

			// test
			if _, err := groups.Datastore.Pool.Exec(t.Context(), "UPDATE "+queue+" SET can_run_after = now() WHERE consumer_group_id = $1 AND message_id = 1", consumer.Id); err != nil {
				t.Fatal(err)
			}
			predecessor, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(predecessor) != 1 || predecessor[0].MessageId != 1 || predecessor[0].Attempts != 1 || predecessor[0].Delays != 1 {
				t.Fatalf("Claim when predecessor is due = %+v, %v; want message 1 at attempt 1, delays 1", predecessor, err)
			}

			// test
			err = exceptions.RecordTerminal(t.Context(), &predecessor[0], errors.New("order cancelled"), mode, nil)

			// verify
			if err != nil {
				t.Fatal(err)
			}

			// test
			second, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(second) != 1 || second[0].MessageId != 2 || second[0].Attempts != 0 || second[0].Delays != 0 {
				t.Fatalf("Claim after terminal predecessor = %+v, %v; want message 2 at attempt 0, delays 0", second, err)
			}

			// test
			err = exceptions.RecordSuperseded(t.Context(), &second[0], mode)

			// verify
			if err != nil {
				t.Fatal(err)
			}

			// test
			third, err := exceptions.Claim(t.Context(), consumer.StreamId, consumer.Id, 1, 100, 3, time.Minute, mode)

			// verify
			if err != nil || len(third) != 1 || third[0].MessageId != 3 || third[0].Attempts != 0 || third[0].Delays != 0 {
				t.Fatalf("Claim after superseded predecessor = %+v, %v; want message 3 at attempt 0, delays 0", third, err)
			}

			// test
			err = exceptions.RecordSuccess(t.Context(), &third[0], mode, nil)

			// verify
			if err != nil {
				t.Fatal(err)
			}
			var history string
			if err := groups.Datastore.Pool.QueryRow(t.Context(), "SELECT coalesce(string_agg(message_id || ':' || status || ':' || attempt, ',' ORDER BY message_id, id), '') FROM "+logs+" WHERE consumer_group_id = $1", consumer.Id).Scan(&history); err != nil {
				t.Fatal(err)
			}
			want := "1:delayed:0,1:failure:1,2:deferred:0,2:superseded:0,3:deferred:0"
			if mode == stream.DeliveryLogModeAll {
				want += ",3:success:0"
			} else if mode == stream.DeliveryLogModeOff {
				want = ""
			}
			if history != want {
				t.Fatalf("ordered history = %q, want %q", history, want)
			}
		})
	}
}
