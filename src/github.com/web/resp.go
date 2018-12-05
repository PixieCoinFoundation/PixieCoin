package web

type restData struct {
	status int
	data   interface{}
}

func NewRestData() *restData {
	return &restData{}
}

func (this *restData) SetStatus(status int) *restData {
	this.status = status
	return this
}

func (this *restData) SetData(data interface{}) *restData {
	this.data = data
	return this
}

func (this *restData) Status() int {
	return this.status
}

func (this *restData) Data() interface{} {
	return this.data
}
