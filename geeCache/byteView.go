package geeCache

type ByteView struct {
	byteArr []byte
}

func (view ByteView) Len() int {
	return len(view.byteArr)
}

// ByteSlice - return a deep copy of the byteArr
func (view ByteView) ByteSlice() []byte {
	myCopy := make([]byte, len(view.byteArr))
	copy(myCopy, view.byteArr)
	return myCopy
}

func (view ByteView) AsString() string {
	return string(view.byteArr)
}
