# Traefik Fault Injection

This plugin can be used to test the resiliency of microservices to different forms of failures.

It's highly inspired by the [Envoy Fault Injection filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/fault_filter).

## Docs

Currently supported header controls:

**x-traefik-fault-delay-request**

The duration to delay a request by. The header value should be an integer that specifies the number of milliseconds to throttle the latency for.

**x-traefik-fault-abort-request**

HTTP status code to abort a request with. The header value should be an integer that specifies the HTTP status code to return in response to a request.

### Plugin options

**Delay**

*Default: true*

This determines if the delay failure is enabled.

**DelayDuration**

*Default: 0*

The number of number of milliseconds to throttle the latency for.

**Abort**

*Default: true*

This determines if the abort failure is enabled.

**AbortCode**

*Default: 400*

The HTTP status code to return.