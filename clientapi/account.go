package clientapi

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/vovamod/go-pterodactyl/api"
	"github.com/vovamod/go-pterodactyl/internal/requester"
)

type accountService struct{ client requester.Requester }

func newAccountService(client requester.Requester) AccountService {
	return &accountService{client: client}
}

func (s *accountService) GetDetails(ctx context.Context) (*api.Account, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/api/client/account", nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.ListItem[api.Account]{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res.Attributes, nil
}

func (s *accountService) GetTwoFactorDetails(ctx context.Context) (*api.TwoFactorDetails, error) {
	req, err := s.client.NewRequest(ctx, "GET", "/api/client/account/two-factor", nil, nil)
	if err != nil {
		return nil, err
	}

	res := &api.TwoFactorDetails{}
	_, err = s.client.Do(ctx, req, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *accountService) EnableTwoFactor(ctx context.Context, options api.TwoFactorEnableOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "POST", "/api/client/account/two-factor", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *accountService) DisableTwoFactor(ctx context.Context, options api.TwoFactorDisableOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "DELETE", "/api/client/account/two-factor", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *accountService) UpdateEmail(ctx context.Context, options api.UpdateEmailOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "PUT", "/api/client/account/email", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *accountService) UpdatePassword(ctx context.Context, options api.UpdatePasswordOptions) error {
	jsonBytes, err := json.Marshal(options)
	if err != nil {
		return err
	}
	req, err := s.client.NewRequest(ctx, "PUT", "/api/client/account/password", bytes.NewBuffer(jsonBytes), nil)
	if err != nil {
		return err
	}
	_, err = s.client.Do(ctx, req, nil)
	return err
}

func (s *accountService) APIKeys() APIKeysService {
	return newAPIKeysService(s.client)
}
