package account

import (
	"context"
	"google.golang.org/grpc/status"
)

// NoFundingCode is the error code returned when a user has no founding. https://grpc.github.io/grpc/core/md_doc_statuscodes.html
const NoFundingCode = 17

func NoFundingError() error {
	return status.Errorf(NoFundingCode, "no funding available, please contact support")
}

func (service *Service) Verify(ctx context.Context) (userId string, err error) {
	return service.Auth.Verify(ctx)
}

func (service *Service) VerifyFunding(ctx context.Context) (userId string, err error) {
	userId, err = service.Auth.Verify(ctx)
	if err != nil {
		return
	}

	balance, err := service.getFinancialSummary(ctx, userId)
	if err != nil {
		return
	}

	if balance.Balance <= 0 {
		err = NoFundingError()
		return
	}

	return
}
