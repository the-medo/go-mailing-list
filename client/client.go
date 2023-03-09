package main

import (
	"context"
	"github.com/alexflint/go-arg"
	pb "go-mailing-list/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func logResponse(res *pb.EmailResponse, err error) {
	if err != nil {
		log.Printf("  error: %v\n", err)
		return
	}

	if res.EmailEntry == nil {
		log.Printf("  no entry\n")
	} else {
		log.Printf("  response: %v\n", res)
	}
}

func createEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("client CreateEmail")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.CreateEmail(ctx, &pb.CreateEmailRequest{EmailAddr: &addr})
	logResponse(res, err)

	return res.EmailEntry
}

func getEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("client GetEmail")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmail(ctx, &pb.GetEmailRequest{EmailAddr: &addr})
	logResponse(res, err)

	return res.EmailEntry
}

func getEmailBatch(client pb.MailingListServiceClient, count int32, page int32) {
	log.Println("client GetEmailBatch")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.GetEmailBatch(ctx, &pb.GetEmailBatchRequest{Count: &count, Page: &page})
	if err != nil {
		log.Fatalf("  error: %v\n", err)
	}
	log.Println("response:")
	for i := 0; i < len(res.EmailEntries); i++ {
		log.Printf(" item [%v of %v]: %s\n", i+1, len(res.EmailEntries), res.EmailEntries[i])
	}
}

func updateEmail(client pb.MailingListServiceClient, entry pb.EmailEntry) *pb.EmailEntry {
	log.Println("client UpdateEmail")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.UpdateEmail(ctx, &pb.UpdateEmailRequest{EmailEntry: &entry})
	logResponse(res, err)

	return res.EmailEntry
}

func deleteEmail(client pb.MailingListServiceClient, addr string) *pb.EmailEntry {
	log.Println("client DeleteEmail")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.DeleteEmail(ctx, &pb.DeleteEmailRequest{EmailAddr: &addr})
	logResponse(res, err)
	return res.EmailEntry
}

var args struct {
	GrpcAddr string `arg:"env:MAILINGLIST_GRPC_ADDR"`
}

func main() {
	arg.MustParse(&args)

	if args.GrpcAddr == "" {
		args.GrpcAddr = ":8081"
	}

	conn, err := grpc.Dial(args.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMailingListServiceClient(conn)

	//confirmedAt := int64(10000)
	//newEmail := createEmail(client, "999@999.997")
	//newEmail.ConfirmedAt = &confirmedAt
	//updateEmail(client, *newEmail)
	//deleteEmail(client, *newEmail.Email)

	getEmailBatch(client, 3, 1)
	getEmailBatch(client, 3, 2)
	getEmailBatch(client, 3, 3)
}
