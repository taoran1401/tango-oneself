package producer

type SimpleProducer struct {
}

func NewSimpleProducer() *SimpleProducer {
	return &SimpleProducer{}
}

func (this *SimpleProducer) Handle() {
	//投递生产
	//global.QUEUE.Publish("simple producer")
}
