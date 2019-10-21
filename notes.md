# Multiply glBufferData size argument
Multiply len(data) * 4 unless it doesnt work. 4 indicates bytes of data type e.g float or int.

# glGen* vs glCreate* and the DSA
The purpose of `direct state access` is not to remove object binding from your application completely (that would be the purpose of the various "bindless" extensions). The purpose of direct state access is to allow you to access the state of an object without having to bind it (ie: directly).

Pre-DSA, you had to bind a buffer just to allocate storage, upload data, etc. With DSA functions, you don't have to. You pass the buffer directly to the function you use to manipulate its state.

But to actually use buffers in a rendering process, you must still bind it to the context or attach it to some other object which itself will get bound to the context.

# Use gl.PtrOffset for something like indices

# Depth_test for not showing back side