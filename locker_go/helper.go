package locker

import (
	"fmt"
)

func prepareKey(prefix string, keyLock KeyLock) (string, error) {
	preparedKey, err := keyLock.GetKeyLockString()
	if err != nil {
		return "", fmt.Errorf("failed to prepare key for %s: %w", keyLock, err)
	}

	return fmt.Sprintf("%s-%s", prefix, preparedKey), nil
}
