Package tattle provides a simple facility to attach an error to the
structure that caused it.

# High level summary

What tattle does could  be done using error values:

	type Record struct {
	    <record contents>
	    err error
	}

Tattlers add the source code location of the error, and a logging
capability.

	type tattler tattle.Tattler // if you like
	...
	type Record struct {
	   <record contents>
	   tat tattler
	}

Prior to an error, a tattler is a nil pointer to a concrete structure - so
it is smaller than even a nil error, out of sensitivity to applications that
have extraordinary object counts.

The tattler's Latch method, tat.Latch(<error>), allocates a bit of memory to
store the error and its source code latch location.  Only the first non-nil
error is latched; subsequent non-nil errors are counted but discarded.

A tattler doesn't implement the error interface. Tattlers are a mousetrap,
not a mouse.  The latched error is available through method Le(). So,
tat.Le().  Yeah, it's a coding pun.

# Experience

Tattlers are in use in a developing database / http server project.
Currently at 14K lines of go code, the project defines 59
tattlers, with 183 Latch cases. The vast majority of Latch cases detect
programmatic, not end-user, errors. The tattlers are helpful.

# Concurrency

Concurrency for tattlers present the same challenge as
concurrency for error variables.  In the above-mentioned project, the
concurrency mechanisms for the structures also protect the tattlers.

Lacking a use case,
I haven't defined a "SerialTattler" type.

# Advice

  - Return errors, not tattlers.
  - The most-used methods are Latch(), Latchf(), Le(), and Logf().
  - A deferred Logf() is your friend.
