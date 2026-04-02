package server

import (
	pb "github.com/ezcnrmn/vaito/gen/go/listing"
	"github.com/ezcnrmn/vaito/services/listing/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func listingToProtobufListing(listing *model.Listing) *pb.Listing {
	pbListing := &pb.Listing{
		Id:          listing.ID,
		Title:       listing.Title,
		Description: listing.Description,
		Category:    &pb.Category{Id: listing.Cetegory.ID, Name: listing.Cetegory.Name},
		UserId:      listing.UserID,
		Status:      listing.Status,
		Price:       listing.Price,
		CreatedAt:   timestamppb.New(listing.CreatedAt),
	}
	if listing.PublishedAt != nil {
		pbListing.PublishedAt = timestamppb.New(*listing.PublishedAt)
	}
	return pbListing
}
