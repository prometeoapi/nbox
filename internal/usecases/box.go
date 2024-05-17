package usecases

import (
	"nbox/internal/domain"
	"nbox/internal/domain/models"
)

type BoxUseCase struct {
	boxOperation domain.StoreOperations
}

func NewBox(boxOperation domain.StoreOperations) *BoxUseCase {
	return &BoxUseCase{boxOperation: boxOperation}
}

func (b *BoxUseCase) Upsert(boxName string, command models.Command[any]) bool {
	//box, _ := b.boxOperation.GetCreateBox(boxName)

	//switch command.Command {
	//case models.UpsertVariable:
	//	return b.UpsertVariable(box, command)
	//case models.UpsertTemplate:
	//	return b.UpsertTemplate(box, command)
	//default:
	//	return false
	//}
	return true
}

func (b *BoxUseCase) UpsertTemplate(box models.Box, command models.Command[any]) bool {
	return true
}

func (b *BoxUseCase) UpsertVariable(box models.Box, command models.Command[any]) bool {
	return true
}
