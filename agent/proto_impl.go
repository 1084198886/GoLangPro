package agent

import (
	"context"
	pb "supwisdom.com/diancan/weigh_dc/proto"
)

func newBindServer(agent *Agent) *bindServerImpl {
	return &bindServerImpl{}
}

type bindServerImpl struct {
	pb.BindDeviceServer
	agent *Agent
}

func newDeviceServer(agent *Agent) pb.WeighScaleDeviceServer {
	return &deviceServiceImpl{
		agent: agent,
	}
}

type deviceServiceImpl struct {
	pb.WeighScaleDeviceServer
	agent *Agent
}

func (d *deviceServiceImpl) Login(context.Context, *pb.WeighingScaleLogin) (*pb.WeighingScaleLoginResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) HeartBeat(context.Context, *pb.WeighingScaleHeartBeat) (*pb.WeighingScaleResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) InitDish(context.Context, *pb.InitDishRequest) (*pb.InitDishResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) ConfirmDish(context.Context, *pb.ConfirmDishRequest) (*pb.ConfirmDishResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) EmitEvent(context.Context, *pb.WeighingScaleEvent) (*pb.WeighingScaleResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) DownloadFood(context.Context, *pb.DownloadFoodRequest) (*pb.FoodResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) DownloadMeal(context.Context, *pb.DownloadMealRequest) (*pb.MealResponse, error) {
	panic("implement me")
}

func (d *deviceServiceImpl) DownloadSysPara(context.Context, *pb.DownloadSysParaRequest) (*pb.SysparaResponse, error) {
	panic("implement me")
}

func (deviceServiceImpl) UploadWeighAction(context.Context, *pb.WeighActionRequest) (*pb.WeighActionResponse, error) {
	panic("implement me")
}
