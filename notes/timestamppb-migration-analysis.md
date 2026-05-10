# Experiment: Protobuf Structs & timestamppb Migration Analysis

## Objectives
- Standardize all time-related fields in the API response layer by replacing `string` types with `google.protobuf.Timestamp`.
- Achieve consistency and minimize formatting errors across the codebase.
- Ensure identical response structures between REST and gRPC interfaces.

## Implementation Strategy
The experiment was conducted on the following handlers:
- **Health & Readiness:** `Healthz` and `Readyz`.
- **Merchant Module:** `Merchant List` and `Merchant Detail`.

To achieve the objectives, the following structural changes were tested:
- Replacing `any` types with `*structpb.Value` in `Status.Data`.
- Refactoring `Response` (and its nested structs) and `Log` in [`pkg/response/`](../pkg/response/) using Protobuf messages.
- Replacing `any` types with `*structpb.Value` in `Response.Data` and `ResponseStatus.Details`.

## Findings

### Pros
- **Consistency:** Time-related fields can be directly assigned using `google.protobuf.Timestamp`, ensuring a unified standard.
- **Efficiency in Casting:** Eliminates the need for manual `time.Time` to `string` (RFC3339) conversions at the handler level.

### Cons
- **Potential Performance Overhead:** Overhead due to reflection, as well as increased complexity in JSON encoding and decoding [`pkg/response/status.go:NewValue(v any)`](../pkg/response/status.go#L231).
- **Swagger/OpenAPI Compatibility:** By default, `timestamppb` is rendered as an object (containing `seconds` and `nanos`) in Swagger, whereas the current API standard requires a flat RFC3339 `string`.
- **Boilerplate:** While it removes string conversion, it still requires explicit conversion from `time.Time` to `timestamppb.New()`.

## Preliminary Conclusion
The current results indicate that the **operational cost and technical trade-offs outweigh the benefits**. 

Specifically, the disruption to JSON serialization and the potential overhead of reflection make this approach less ideal for the current state of the "Resik" project. This experiment will be deferred until a more efficient way to handle JSON marshaling for Protobuf timestamps is identified or until the project fully pivots toward a gRPC-first architecture.

---
*This document is preserved for future reference and documentation of the R&D process.*
