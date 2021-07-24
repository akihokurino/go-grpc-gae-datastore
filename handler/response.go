package handler

import (
	"net/url"

	"gae-go-sample/domain"
	pb "gae-go-sample/proto/go/pb"
)

func urlString(url *url.URL, alt string) string {
	if url == nil || url.String() == "" {
		return alt
	}

	return url.String()
}

func ToApplyClientResponse(from *domain.ApplyClient) *pb.ApplyClient {
	return &pb.ApplyClient{
		Email:           from.Email.String(),
		PhoneNumber:     from.PhoneNumber,
		CompanyName:     from.CompanyName,
		WebURL:          from.WebURL.String(),
		AccountName:     from.AccountName,
		AccountNameKana: from.AccountNameKana,
		Position:        from.Position,
		Status:          from.Status,
		CreatedAt:       domain.ToJSTString(from.CreatedAt),
	}
}

func ToCompanyResponse(from *domain.Company) *pb.Company {
	return &pb.Company{
		Id:                         from.ID.String(),
		Name:                       from.Name,
		GsLogoURL:                  urlString(from.GSLogoURL, ""),
		SignedLogoURL:              urlString(from.SignedLogoURL, ""),
		WebURL:                     urlString(from.WebURL, ""),
		EstablishedAt:              domain.ToJSTString(from.EstablishedAt),
		PostalCode:                 from.PostalCode,
		RepresentativeName:         from.RepresentativeName,
		CapitalStock:               from.CapitalStock,
		Introduction:               from.Introduction,
		AccordingCompanyName:       from.AccordingCompanyName,
		AccordingCompanyPostalCode: from.AccordingCompanyPostalCode,
		AccordingCompanyAddress:    from.AccordingCompanyAddress,
		Status:                     from.Status,
		CreatedAt:                  domain.ToJSTString(from.CreatedAt),
	}
}

func ToProjectResponse(from *domain.Project) *pb.Project {
	highlights := make([]*pb.SearchHighlight, 0, len(from.Highlights))
	for _, highlight := range from.Highlights {
		highlights = append(highlights, &pb.SearchHighlight{
			Key:       highlight.Key,
			Val:       highlight.Val,
			MatchWord: highlight.MatchedWord,
			LinkId:    highlight.LinkID,
		})
	}

	return &pb.Project{
		Id:                 from.ID.String(),
		CompanyID:          from.CompanyID.String(),
		Name:               from.Name,
		Description:        from.Description,
		GsThumbnailURL:     urlString(from.GSThumbnailURL, ""),
		SignedThumbnailURL: urlString(from.SignedThumbnailURL, ""),
		Status:             from.Status,
		CreatedAt:          domain.ToJSTString(from.CreatedAt),
		Highlights:         highlights,
	}
}

func ToClientResponse(from *domain.Client) *pb.Client {
	return &pb.Client{
		Id:            from.ID.String(),
		CompanyID:     from.CompanyID.String(),
		Name:          from.Name,
		NameKana:      from.NameKana,
		GsIconURL:     urlString(from.GSIconURL, ""),
		SignedIconURL: urlString(from.SignedIconURL, ""),
		PhoneNumber:   from.PhoneNumber,
		Position:      from.Position,
		Email:         from.Email,
		Role:          from.Role,
		CreatedAt:     domain.ToJSTString(from.CreatedAt),
	}
}

func ToCustomerResponse(from *domain.Customer) *pb.Customer {
	highlights := make([]*pb.SearchHighlight, 0, len(from.Highlights))
	for _, highlight := range from.Highlights {
		highlights = append(highlights, &pb.SearchHighlight{
			Key:       highlight.Key,
			Val:       highlight.Val,
			MatchWord: highlight.MatchedWord,
			LinkId:    highlight.LinkID,
		})
	}

	return &pb.Customer{
		Id:            from.ID.String(),
		Name:          from.Name,
		NameKana:      from.NameKana,
		GsIconURL:     urlString(from.GSIconURL, ""),
		SignedIconURL: urlString(from.SignedIconURL, ""),
		Birthdate:     domain.ToJSTString(from.Birthdate),
		Gender:        from.Gender,
		PhoneNumber:   from.PhoneNumber,
		Email:         from.Email,
		Pr:            from.Pr,
		Address:       from.Address,
		Status:        from.Status,
		ResumeURL:     urlString(from.ResumeURL, ""),
		CreatedAt:     domain.ToJSTString(from.CreatedAt),
		Highlights:    highlights,
	}
}

func ToEntryResponse(from *domain.Entry) *pb.Entry {
	return &pb.Entry{
		CustomerID: from.CustomerID.String(),
		ProjectID:  from.ProjectID.String(),
		CreatedAt:  domain.ToJSTString(from.CreatedAt),
	}
}

func ToMessageRoomResponse(from *domain.MessageRoom) *pb.MessageRoom {
	return &pb.MessageRoom{
		Id:                  from.ID().String(),
		ProjectID:           from.ProjectID.String(),
		CustomerID:          from.CustomerID.String(),
		CompanyID:           from.CompanyID.String(),
		CustomerUnRead:      from.CustomerUnRead,
		CustomerUnReadCount: from.CustomerUnReadCount,
		CompanyUnRead:       from.CompanyUnRead,
		CompanyUnReadCount:  from.CompanyUnReadCount,
		CreatedAt:           domain.ToJSTString(from.CreatedAt),
		UpdatedAt:           domain.ToJSTString(from.UpdatedAt),
	}
}

func ToMessageResponse(from *domain.Message) *pb.Message {
	return &pb.Message{
		Id:             from.ID.String(),
		RoomID:         from.RoomID.String(),
		FromID:         from.FromID,
		ToID:           from.ToID,
		FromCompany:    from.FromCompany,
		Text:           from.Text,
		GsImageURL:     urlString(from.GSImageURL, ""),
		SignedImageURL: urlString(from.SignedImageURL, ""),
		GsFileURL:      urlString(from.GSFileURL, ""),
		SignedFileURL:  urlString(from.SignedFileURL, ""),
		CreatedAt:      domain.ToJSTString(from.CreatedAt),
	}
}

func ToPrefectureResponse(from *domain.Prefecture) *pb.Prefecture {
	return &pb.Prefecture{
		Code: from.Code,
		Name: from.Name,
	}
}

func ToContractResponse(from *domain.Contract) *pb.Contract {
	return &pb.Contract{
		CompanyID:     from.CompanyID.String(),
		ProjectID:     from.ProjectID.String(),
		CustomerID:    from.CustomerID.String(),
		GsFileURL:     urlString(from.GSFileURL, ""),
		SignedFileURL: urlString(from.SignedFileURL, ""),
		Status:        from.Status,
		CreatedAt:     domain.ToJSTString(from.CreatedAt),
		UpdatedAt:     domain.ToJSTString(from.UpdatedAt),
	}
}
