# Layer 8 Parser (L8Parser)

A model-agnostic parsing engine for the Layer 8 ecosystem. L8Parser processes raw data collected by [l8collector](https://github.com/saichler/l8collector) — SNMP walks, SSH command outputs, REST API responses, and kubectl queries — and transforms it into structured protobuf objects (`NetworkDevice`, `GpuDevice`, `K8sCluster`) using a configurable rule-based framework.

## Architecture

### Processing Flow

1. **Job Reception** — Parser receives completed collection jobs from l8collector
2. **Rule Application** — Applies parsing rules based on [l8pollaris](https://github.com/saichler/l8pollaris) configuration
3. **Data Transformation** — Transforms raw data into structured protobuf objects via PropertyId injection
4. **Inventory Update** — Sends parsed objects to [l8inventory](https://github.com/saichler/l8inventory) via PATCH operations

### Core Components

| Component | Path | Purpose |
|-----------|------|---------|
| Parser | `go/parser/service/Parser.go` | Main parsing engine that processes jobs |
| ParsingService | `go/parser/service/ParsingService.go` | Layer 8 service interface wrapper |
| ParsingCenter | `go/parser/service/ParsingCenter.go` | Job completion handler and inventory integration |
| Rule Engine | `go/parser/rules/` | 18 parsing rule implementations |
| Boot Configs | `go/parser/boot/` | 21 vendor-specific polling configurations |

### Project Structure

```
l8parser/
├── go/
│   ├── parser/
│   │   ├── boot/                        # Vendor-specific polling configs
│   │   │   ├── SNMP.go                  # Common SNMP utilities, boot stages, cadence plans
│   │   │   ├── equipmentinfo.go         # Equipment info attribute helpers
│   │   │   ├── K8s.go                   # Kubernetes resource monitoring
│   │   │   ├── nvidia.go               # NVIDIA GPU SNMP polling
│   │   │   ├── nvidia_ssh_rest.go      # NVIDIA GPU SSH + REST polling
│   │   │   ├── ospf_bgp_vrf.go         # OSPF/BGP/VRF polling definitions
│   │   │   ├── cisco.go                # Cisco switches and routers
│   │   │   ├── juniper.go              # Juniper routers
│   │   │   ├── paloalto.go             # Palo Alto firewalls
│   │   │   ├── fortinet.go             # Fortinet firewalls
│   │   │   ├── arista.go               # Arista switches
│   │   │   ├── nokia.go                # Nokia routers
│   │   │   ├── huawei.go               # Huawei routers
│   │   │   ├── dell.go                 # Dell servers
│   │   │   ├── hpe.go                  # HPE servers
│   │   │   ├── ibm.go                  # IBM servers
│   │   │   ├── checkpoint.go           # Check Point firewalls
│   │   │   ├── sonicwall.go            # SonicWall firewalls
│   │   │   ├── dlink.go                # D-Link switches
│   │   │   ├── extreme.go              # Extreme Networks switches
│   │   │   └── nec.go                  # NEC routers
│   │   ├── rules/                       # Parsing rule implementations
│   │   │   ├── ParsingRule.go           # Rule interface and registry
│   │   │   ├── ParamNames.go            # Shared parameter name constants
│   │   │   ├── Contains.go             # Text pattern matching
│   │   │   ├── Set.go                  # Direct value assignment
│   │   │   ├── NormalizeEnum.go        # Enum value normalization
│   │   │   ├── StringToCTable.go       # String to columnar table conversion
│   │   │   ├── CTableToMapProperty.go  # Table to map property transform
│   │   │   ├── EntityMibToPhysicals.go # SNMP Entity MIB parsing
│   │   │   ├── IfTableToPhysicals.go   # SNMP ifTable parsing
│   │   │   ├── SnmpGpuTable.go         # SNMP GPU table parsing
│   │   │   ├── SnmpOspfToVrf.go        # SNMP OSPF MIB to VRF parsing
│   │   │   ├── SnmpBgpToVrf.go         # SNMP BGP MIB to VRF parsing
│   │   │   ├── SshNvidiaSmiParse.go    # nvidia-smi SSH output parsing
│   │   │   ├── SshVrfParse.go          # Multi-vendor "show vrf" SSH parsing
│   │   │   ├── RestJsonParse.go        # Generic REST JSON response parsing
│   │   │   ├── RestGpuParse.go         # GPU REST API parsing (DCGM)
│   │   │   ├── InferDeviceType.go      # Device type inference from sysOID
│   │   │   ├── MapToDeviceStatus.go    # Device status mapping
│   │   │   └── SetTimeSeries.go        # Time-series metric handling
│   │   └── service/                     # Core parsing services
│   │       ├── Parser.go
│   │       ├── ParsingService.go
│   │       └── ParsingCenter.go
│   ├── tests/                           # Test suite
│   │   ├── jobsPersistency/            # Persistent real device data for replay tests
│   │   ├── TestInit.go                 # Test initialization and topology setup
│   │   ├── Parser_test.go
│   │   ├── DevicesParsing_test.go
│   │   ├── PhysicalTest_test.go
│   │   ├── PhysicalTestFromPersistency_test.go
│   │   ├── Property_test.go
│   │   ├── TestDevices_test.go
│   │   ├── ClusterTest_test.go
│   │   └── Devices.go
│   ├── go.mod
│   ├── go.sum
│   └── test.sh                          # Build, test, and coverage script
├── networkdevice-unpopulated-attributes.md
├── vendor-specific-oids.md
└── README.md
```

## Parsing Rules

L8Parser provides 18 parsing rules covering four collection protocols:

### Generic Rules
| Rule | Purpose |
|------|---------|
| Contains | Matches text patterns in collected data |
| Set | Directly assigns values to object properties |
| NormalizeEnum | Normalizes raw values to enum constants |
| MapToDeviceStatus | Maps protocol-specific status codes to device status |
| SetTimeSeries | Handles time-series metric injection |
| InferDeviceType | Infers device type from sysOID enterprise prefix |

### SNMP Rules
| Rule | Purpose |
|------|---------|
| StringToCTable | Converts SNMP table walks to columnar tables |
| CTableToMapProperty | Transforms columnar tables to map properties |
| EntityMibToPhysicals | Parses SNMP Entity MIB into physical components |
| IfTableToPhysicals | Parses SNMP ifTable into logical interfaces |
| SnmpGpuTable | Parses SNMP GPU tables (NVIDIA enterprise MIB) |
| SnmpOspfToVrf | Parses OSPF MIB into VRF structures |
| SnmpBgpToVrf | Parses BGP MIB into VRF structures |

### SSH Rules
| Rule | Purpose |
|------|---------|
| SshNvidiaSmiParse | Parses `nvidia-smi` command output (utilization, temperature, power) |
| SshVrfParse | Parses `show vrf` output across 9 vendor-specific formats |

### REST Rules
| Rule | Purpose |
|------|---------|
| RestJsonParse | Generic REST JSON response parser |
| RestGpuParse | NVIDIA DCGM REST API parser (topology, health, NVLink) |

## Vendor Support

### Network Devices (SNMP + SSH)

15 vendors with automatic detection by sysOID enterprise prefix:

| Vendor | Type | Enterprise OID |
|--------|------|----------------|
| Cisco | Switches, Routers | `.1.3.6.1.4.1.9.` |
| Juniper | Routers | `.1.3.6.1.4.1.2636.` |
| Palo Alto | Firewalls | `.1.3.6.1.4.1.25461.` |
| Fortinet | Firewalls | `.1.3.6.1.4.1.12356.` |
| Arista | Switches | `.1.3.6.1.4.1.30065.` |
| Nokia | Routers | `.1.3.6.1.4.1.6527.` |
| Huawei | Routers | `.1.3.6.1.4.1.2011.` |
| Dell | Servers | `.1.3.6.1.4.1.674.` |
| HPE | Servers | `.1.3.6.1.4.1.232.` |
| IBM | Servers | `.1.3.6.1.4.1.2.` |
| Check Point | Firewalls | `.1.3.6.1.4.1.2620.` |
| SonicWall | Firewalls | `.1.3.6.1.4.1.8741.` |
| D-Link | Switches | `.1.3.6.1.4.1.171.` |
| Extreme | Switches | `.1.3.6.1.4.1.1916.` |
| NEC | Routers | `.1.3.6.1.4.1.119.` |

Each vendor configuration polls for:
- System information (vendor, version, serial numbers)
- Interface monitoring (status, speed, MTU, names)
- Hardware inventory (modules, power supplies, fans, temperature)
- Performance metrics (CPU, memory, utilization)

### NVIDIA GPU (Three-Protocol Stack)

Full GPU monitoring via SNMP, SSH, and REST:

| Protocol | Polls | Data Collected |
|----------|-------|----------------|
| SNMP | 7 | System info, GPU modules, per-GPU static info, metrics, host resources, interfaces, status |
| SSH | 5 | nvidia-smi (utilization, temperature, power), show version, lscpu |
| REST | 4 | DCGM device info, topology/NVLink, health, memory |

### Kubernetes (kubectl)

Resource monitoring via kubectl commands:

| Resource | Command |
|----------|---------|
| Nodes | `get nodes -o wide` |
| Pods | `get pods -A -o wide` |
| Deployments | `get deployments -A -o wide` |
| StatefulSets | `get statefulsets -A -o wide` |
| DaemonSets | `get daemonsets -A -o wide` |
| Services | `get services -A -o wide` |
| Namespaces | `get namespaces` |
| Network Policies | `get networkpolicies -A -o wide` |

Features include log collection, on-demand detail queries, and multi-cluster support.

### OSPF/BGP/VRF

Routing protocol support across vendors:

- **SNMP**: Standard MIB-II OIDs for OSPF neighbor and BGP peer discovery (vendor-agnostic)
- **SSH**: `show vrf` parsing for 9 vendor-specific output formats (Cisco IOS-XR/IOS-XE/NX-OS, Juniper, Nokia, Huawei, Arista, Extreme, NEC)

## Usage

### Basic Integration

```go
import (
    "github.com/saichler/l8parser/go/parser/service"
)

// Create and activate parsing service
parsingService := &service.ParsingService{}
err := parsingService.Activate(
    "parsing-service",
    serviceArea,
    resources,
    vnic,
    &NetworkDevice{},
    "Id",
)
```

### Polling Configuration

```go
poll := &types.Poll{
    Name:      "systemMib",
    What:      ".1.3.6.1.2.1.1",
    Operation: types.Operation_OMap,
    Attributes: []*types.Attribute{
        {
            PropertyId: "networkdevice.equipmentinfo.vendor",
            Rules: []*types.Rule{
                {
                    Name: "Contains",
                    Params: map[string]*types.Parameter{
                        "what":   {Value: "cisco"},
                        "from":   {Value: ".1.3.6.1.2.1.1.1.0"},
                        "output": {Value: "Cisco"},
                    },
                },
            },
        },
    },
}
```

### Device Detection

```go
// Automatic vendor detection by sysOID
polaris := GetPollarisByOid(".1.3.6.1.4.1.9.1.122.0") // Cisco switch config

// Get all vendor models
models := GetAllPolarisModels()
```

## Prerequisites

- Go 1.25+ (current: Go 1.25.4)
- Access to Layer 8 ecosystem modules

### Dependencies

| Module | Purpose |
|--------|---------|
| [l8collector](https://github.com/saichler/l8collector) | Data collection framework |
| [l8inventory](https://github.com/saichler/l8inventory) | Inventory management |
| [l8pollaris](https://github.com/saichler/l8pollaris) | Polling configuration |
| [l8types](https://github.com/saichler/l8types) | Common types and interfaces |
| [l8bus](https://github.com/saichler/l8bus) | Message bus |
| [l8utils](https://github.com/saichler/l8utils) | Shared utilities |
| [l8reflect](https://github.com/saichler/l8reflect) | Reflection utilities |
| [l8srlz](https://github.com/saichler/l8srlz) | Serialization |
| [probler](https://github.com/saichler/probler) | Network probing |

## Build & Test

```bash
cd go
./test.sh
```

The script initializes Go modules, fetches dependencies, runs the full test suite with coverage, and opens a coverage report.

### Test Suite

Tests use persistent real-device data from `go/tests/jobsPersistency/` for replay-based validation:

- **Parser_test.go** — Core parser engine tests
- **DevicesParsing_test.go** — Vendor-specific device parsing
- **PhysicalTest_test.go** — Physical component extraction
- **PhysicalTestFromPersistency_test.go** — Real device data replay
- **Property_test.go** — PropertyId injection
- **TestDevices_test.go** — Device type inference
- **ClusterTest_test.go** — Kubernetes cluster parsing

## License

© 2025-2026 Sharon Aicler (saichler@gmail.com)

Licensed under the Apache License, Version 2.0. See [LICENSE](http://www.apache.org/licenses/LICENSE-2.0).
