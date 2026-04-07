package locker

import (
	"asiatix/internal/metrics"
	"context"
)

type Wrapper struct {
	locker   Locker
	exporter metrics.Exporter
}

func newLockerWrapper(locker Locker, exporter metrics.Exporter) Locker {
	return &Wrapper{locker, exporter}
}

func (w *Wrapper) Lock(keyLock KeyLock) error {
	if err := w.locker.Lock(keyLock); err != nil {
		w.exporter.UpCounterInc(context.TODO(), lockerErrorMetric)

		return err
	}
	return nil
}

func (w *Wrapper) Unlock(keyLock KeyLock) error {
	if err := w.locker.Unlock(keyLock); err != nil {
		w.exporter.UpCounterInc(context.TODO(), lockerErrorMetric)

		return err
	}
	return nil
}
