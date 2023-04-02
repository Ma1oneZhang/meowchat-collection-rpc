package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/scheduled"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCatLogic {
	return &CreateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateCatLogic) CreateCat(in *pb.CreateCatReq) (*pb.CreateCatResp, error) {
	cat := &model.Cat{}
	err := copier.Copy(cat, in.Cat)
	if err != nil {
		return nil, err
	}
	err = l.svcCtx.CatModel.Insert(l.ctx, cat)
	if err != nil {
		return nil, err
	}
	go scheduled.SendUrlUsedMessageToSts(&l.svcCtx.Config, &cat.Avatars)
	return &pb.CreateCatResp{CatId: cat.ID.Hex()}, nil
}
