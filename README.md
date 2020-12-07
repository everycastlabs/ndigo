Go bindings for the [NDI SDK](https://ndi.tv/sdk/) by NewTek

# Getting started
In order to use these bindings the NDI SDK must first be installed for your system from NewTek directly.

This should appear on the path as `ndi.4` and should be discoverable by standard build tools.

To check this is working, try a couple of the examples in the `examples` folder.

# Known issues
## Video frame
* `VideoFrameT` has been deprecated, please avoid using it
* `VideoFrameV2T` is its replacement, but the `(C) line_stride_in_bytes` attribute commonly used in the examples is not
available for Go to bind to due to being in a union.
* Some examples such as `Send_Latency` update the frame attributes in place to reduce overhead and memory allocations.
If this is done, `frame.Free()` must be called before the frame is reused else copied attributes will not be updated.
