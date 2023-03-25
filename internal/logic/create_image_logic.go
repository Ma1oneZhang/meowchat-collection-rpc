package logic

import (
	"context"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateImageLogic {
	return &CreateImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateImageLogic) CreateImage(in *pb.CreateImageReq) (*pb.CreateImageResp, error) {
	data := make([]*model.Image, len(in.Images))
	for i := 0; i < len(data); i++ {
		data[i] = &model.Image{
			CatId:    in.Images[i].CatId,
			ImageUrl: in.Images[i].Url,
		}
	}
	err := l.svcCtx.ImageModel.InsertMany(l.ctx, data)
	if err != nil {
		return nil, err
	}
	id := make([]string, len(data))
	for i := 0; i < len(data); i++ {
		id[i] = data[i].ID.Hex()
		// 将其加入已经使用的图片集合之中
		addToUsedUrl(&l.svcCtx.Config.Redis, data[i].ImageUrl)
	}
	return &pb.CreateImageResp{ImageIds: id}, nil
}
