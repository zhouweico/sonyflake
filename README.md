Sonyflake
=========

[![GoDoc](https://godoc.org/github.com/sony/sonyflake?status.svg)](http://godoc.org/github.com/sony/sonyflake)
[![Build Status](https://travis-ci.org/sony/sonyflake.svg?branch=master)](https://travis-ci.org/sony/sonyflake)
[![Coverage Status](https://coveralls.io/repos/sony/sonyflake/badge.svg?branch=master&service=github)](https://coveralls.io/github/sony/sonyflake?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/sony/sonyflake)](https://goreportcard.com/report/github.com/sony/sonyflake)

Sonyflake is a distributed unique ID generator inspired by [Twitter's Snowflake](https://blog.twitter.com/2010/announcing-snowflake).  
A Sonyflake ID is composed of

    41 bits for time in units of 10 msec
     6 bits for a sequence number
    16 bits for a machine id


sonyflake 中文描述

sonyflake的结构如下(每部分用-分开):

```
+--------------------------------------------------------------------------+
| 1 Bit Unused | 41 Bit Timestamp |  6 Bit Sequence ID  |   16 Bit Node ID |
+--------------------------------------------------------------------------+
```

 * 1位标识，最高位是符号位，正数是0，负数是1，所以id一般是正数，最高位是0
 * 41位时间截(毫秒级)，注意，41位时间截不是存储当前时间的时间截，而是存储时间截的差值（当前时间截 - 开始时间截)
 * 得到的值），这里的的开始时间截，一般是我们的id生成器开始使用的时间，由我们程序来指定的。41位的时间截，可以使用69年，年T = (1L << 41) / (1000L * 60 * 60 * 24 * 365) = 69
 * 6位序列，毫秒内的计数，6位的计数顺序号支持每个节点每毫秒(同一机器，同一时间截)产生64个ID序号
 * 16位的数据机器位，可以部署在65535个节点，网段及子网10.0.0.0/16，172.16-31.0.0/16，192.168.0.0/16
 * 加起来刚好64位，为一个Long型。
 * sonyflake的优点是，整体上按照时间自增排序，单个集群内不会产生ID碰撞。
 
 
Installation
------------

```
go get github.com/zhouweico/sonyflake
```

Usage
-----

The function NewSonyflake creates a new Sonyflake instance.

```go
func NewSonyflake(st Settings) *Sonyflake
```

You can configure Sonyflake by the struct Settings:

```go
type Settings struct {
	StartTime      time.Time
	MachineID      func() (uint16, error)
	CheckMachineID func(uint16) bool
}
```

- StartTime is the time since which the Sonyflake time is defined as the elapsed time.
  If StartTime is 0, the start time of the Sonyflake is set to "2019-1-1 00:00:00 +0000 UTC".
  If StartTime is ahead of the current time, Sonyflake is not created.

- MachineID returns the unique ID of the Sonyflake instance.
  If MachineID returns an error, Sonyflake is not created.
  If MachineID is nil, default MachineID is used.
  Default MachineID returns the lower 16 bits of the private IP address.

- CheckMachineID validates the uniqueness of the machine ID.
  If CheckMachineID returns false, Sonyflake is not created.
  If CheckMachineID is nil, no validation is done.

In order to get a new unique ID, you just have to call the method NextID.

```go
func (sf *Sonyflake) NextID() (uint64, error)
```

NextID can continue to generate IDs for about 174 years from StartTime.
But after the Sonyflake time is over the limit, NextID returns an error.

AWS VPC and Docker
------------------

The [awsutil](https://github.com/zhouweico/sonyflake/blob/master/awsutil) package provides
the function AmazonEC2MachineID that returns the lower 16-bit private IP address of the Amazon EC2 instance.
It also works correctly on Docker
by retrieving [instance metadata](http://docs.aws.amazon.com/en_us/AWSEC2/latest/UserGuide/ec2-instance-metadata.html).

[AWS VPC](http://docs.aws.amazon.com/en_us/AmazonVPC/latest/UserGuide/VPC_Subnets.html)
is assigned a single CIDR with a netmask between /28 and /16.
So if each EC2 instance has a unique private IP address in AWS VPC,
the lower 16 bits of the address is also unique.
In this common case, you can use AmazonEC2MachineID as Settings.MachineID.

See [example](https://github.com/zhouweico/sonyflake/blob/master/example) that runs Sonyflake on AWS Elastic Beanstalk.

License
-------

The MIT License (MIT)

See [LICENSE](https://github.com/zhouweico/sonyflake/blob/master/LICENSE) for details.
