package utils

import "io"

func SafeClose(closer io.Closer) {
	if closer != nil {
		if err := closer.Close(); err != nil {
			Logger().Error("fail to close", ErrLog(err))
		}
	}
}
