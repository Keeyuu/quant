package util

import (
	"app/infrastructure/config"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/satori/go.uuid"
	"sync"
)

var awsSession *session.Session
var awsSessionOnce sync.Once

// 默认使用机器授权
func getAwsSession() (*session.Session, error) {
	var err error
	region := config.Get().Sqs.Region
	awsSessionOnce.Do(func() {
		awsSession, err = session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				Region: &region,
			},
			SharedConfigState: session.SharedConfigEnable,
		})
	})
	return awsSession, err
}

// pull message from aws-sqs.
// queueUrl - The URL of the Amazon SQS queue from which messages are received.
// msgCount - The maximum number of messages to return.
func ReceiveSqsMessage(queueUrl string, msgCount int64) (messages []*sqs.Message, err error) {
	sess, err := getAwsSession()
	if err != nil {
		err = errors.New("download data, create aws session error: " + err.Error())
		return
	}
	svc := sqs.New(sess)
	receiveResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		MaxNumberOfMessages: aws.Int64(msgCount),
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:          aws.String(queueUrl),
		VisibilityTimeout: aws.Int64(60),
		WaitTimeSeconds:   aws.Int64(20),
	})
	if err != nil {
		return
	}
	messages = receiveResult.Messages
	return
}

func DeleteSqsMessage(queueUrl string, message *sqs.Message) (err error) {
	sess, err := getAwsSession()
	if err != nil {
		err = errors.New("download data, create aws session error: " + err.Error())
		return
	}
	svc := sqs.New(sess)
	_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	return
}

func DeleteSqsMessageBatch(queueUrl string, messages []*sqs.Message) (failedItems []*sqs.BatchResultErrorEntry, err error) {
	sess, err := getAwsSession()
	if err != nil {
		err = errors.New("download data, create aws session error: " + err.Error())
		return
	}
	entries := make([]*sqs.DeleteMessageBatchRequestEntry, 0, len(messages))
	for k := range messages {
		entries = append(entries, &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(uuid.NewV4().String()),
			ReceiptHandle: messages[k].ReceiptHandle,
		})
	}
	svc := sqs.New(sess)
	deleteResult, err := svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(queueUrl),
	})
	failedItems = deleteResult.Failed
	return
}
