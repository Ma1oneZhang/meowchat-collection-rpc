package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/xh-polaris/meowchat-collection-rpc/errorx"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/model"
	"github.com/xh-polaris/meowchat-collection-rpc/internal/scheduled"
	"github.com/xh-polaris/meowchat-collection-rpc/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/xh-polaris/meowchat-collection-rpc/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCatLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCatLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCatLogic {
	return &UpdateCatLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCatLogic) UpdateCat(in *pb.UpdateCatReq) (*pb.UpdateCatResp, error) {
	cat := &model.Cat{}
	err := copier.Copy(cat, in.Cat)
	if err != nil {
		return nil, err
	}
	cat.ID, err = primitive.ObjectIDFromHex(in.Cat.Id)
	if err != nil {
		return nil, errorx.ErrInvalidId
	}
	err = l.svcCtx.CatModel.Update(l.ctx, cat)
	if err != nil {
		return nil, err
	}
	go scheduled.SendUrlUsedMessageToSts(&l.svcCtx.Config, &cat.Avatars)
	return &pb.UpdateCatResp{}, nil
}
