package cmd

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
)

func getTime(original string) (parsed time.Time) {
	parsed, _ = time.Parse(
		RFC8601,
		original,
	)
	return
}

var exampleSnapshot1Type = aws.String("db.m1.medium")
var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")
var exampleSnapshot1Time = "2015-06-11T22:00:00+00:00"
var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot1DBID,
	DBSnapshotIdentifier: exampleSnapshot1ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot1Time)),
	Status:               aws.String("available"),
}

var exampleSnapshot2DBID = aws.String("develop-rds")
var exampleSnapshot2ID = aws.String("rds:develop-2016-06-11")
var exampleSnapshot2Time = "2016-06-11T22:00:00+00:00"
var exampleSnapshot2 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot2DBID,
	DBSnapshotIdentifier: exampleSnapshot2ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot2Time)),
	Status:               aws.String("available"),
}

var exampleSnapshot3DBID = aws.String("develop-rds")
var exampleSnapshot3ID = aws.String("rds:develop-2017-06-11")
var exampleSnapshot3Time = "2017-06-11T22:00:00+00:00"
var exampleSnapshot3 = &rds.DBSnapshot{
	DBInstanceIdentifier: exampleSnapshot3DBID,
	DBSnapshotIdentifier: exampleSnapshot3ID,
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot3Time)),
}

var exampleSnapshot4DBID = aws.String("develop-rds")
var exampleSnapshot4ID = aws.String("rds:develop-2017-07-11")
var exampleSnapshot4 = &rds.DBSnapshot{
	DBInstanceIdentifier: exampleSnapshot4DBID,
	DBSnapshotIdentifier: exampleSnapshot4ID,
}
