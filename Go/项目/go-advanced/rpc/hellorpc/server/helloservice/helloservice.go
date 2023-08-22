package helloservice

type HelloService struct{}

func (s HelloService) Hello(request string, reply *string) error {
	*reply = "hello" + request
	return nil
}
