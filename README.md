# rhsm-cli

# Getting a token
Please see [Generating an offline token](https://access.redhat.com/management/api) to generate a token.

```bash
export RHSM_TOKEN=<your offline token>
```

# Building
I'll set up goreleaser later, but for now just do a go build
```bash
go build .
```

# Using rhsm-cli

## Listing all systems under account
```bash
./rhsm-cli list systems
```

## List systems matching a filter
```bash
./rhsm-cli list systems --filter ocp
```

## Show system details by UUID
**Note:** There were some systems that had duplicate hostnames, by using the system UUID instead hostname, it ensures that it brings back the system you intended.

```bash
./rhsm-cli list systems --systemID <SYSTEM UUID>
Hostname, UUID, Subscription Name, Sku, Pool ID
ocp-app-01l.lab1.example.com, SYSTEM_UUID, Red Hat Developer Subscription for Individuals, RH00798, POOL_ID
```

## List subscriptions under account
```bash
./rhsm-cli list subscriptions
Name, Subscription Number, SKU, Status, Pool ID, Quantity, Consumed
Red Hat Enterprise Linux Developer Suite, 4534066, RH2262474, Active
Red Hat Enterprise Linux Developer Suite, 4799758, RH2262474, Active
Red Hat Developer Subscription for Individuals, 9114433, RH00798, Active, 2c928081790fb14c01792457f16a2741, 16, 16
Red Hat Developer Subscription for Individuals, 9122320, RH00798, Active, 8a85f9997922d86501793771efe71fd2, 16, 16
Red Hat Developer Subscription for Individuals, 9340692, RH00798, Active, 8a85f99a7aaf8439017ace705f9b0fbf, 16, 16
Red Hat Enterprise Linux Server, Premium (Physical or Virtual Nodes), 9479832, RH00003S, Active, 8a85f9a07b54e268017be4927ec85f91, 4, 4
```

## Searching subscriptions under account
```bash
$ ./rhsm-cli list subscriptions | grep OpenShift
Red Hat OpenShift Container Platform Standard (2 Cores or 4 vCPUs), 10391081, MCT2736S, Active, 8a85f9997d484aeb017d6d59133e12fd, 46, 46
Red Hat OpenShift Container Platform Premium (2 Cores or 4 vCPUs), 10391080, MCT2735S, Active, 8a85f9997d484aeb017d6d591e031303, 71, 0
Red Hat OpenShift Container Platform Standard (2 Cores or 4 vCPUs), 10391083, MCT2736S, Active, 8a85f9997d484aeb017d6d591ebd1305, 28, 28
Red Hat OpenShift Container Platform Standard (2 Cores or 4 vCPUs), 10391086, MCT2736S, Active, 8a85f9997d484aeb017d6d591a6f1301, 62, 62
Red Hat OpenShift Container Platform Premium (2 Cores or 4 vCPUs), 10391089, MCT2735S, Active, 8a85f9997d484aeb017d6d5939591309, 132, 72
Red Hat OpenShift Container Platform Premium (2 Cores or 4 vCPUs), 10391108, MCT2735S, Active, 8a85f9997d484aeb017d6d58dc6712e8, 71, 0
```

