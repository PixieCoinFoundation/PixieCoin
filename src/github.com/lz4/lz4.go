package lz4

import (
	"fmt"
	"unsafe"
	//    "fmt"
)

const MEMORYUSAGE = 16
const HASHTABLESIZE = 1 << (MEMORYUSAGE - 2)
const HASHLOG = MEMORYUSAGE - 2
const GOLDNUMBER = 2654435761
const LASTLITERALS = 5
const MINMATCH = 4
const COPYLENGTH = 8
const MFLIMIT = COPYLENGTH + MINMATCH
const LZ4MINLENGTH = MFLIMIT + 1

// forward from the start, even if the dst and src overlap.
// It is equivalent to:
//   for i := 0; i < n; i++ {
//     mem[dst+i] = mem[src+i]
//   }
func overlapCopy(mem []byte, dst, src, n uint32) {
	// in LZ4, src is always < dst
	for {
		if dst >= src+n {
			copy(mem[dst:dst+n], mem[src:src+n])
			return
		}
		// There is some forward overlap.  The destination
		// will be filled with a repeated pattern of mem[src:src+k].
		// We copy one instance of the pattern here, then repeat.
		// Each time around this loop k will double.
		k := dst - src
		copy(mem[dst:dst+k], mem[src:src+k])
		n -= k
		dst += k
	}
}
func DecodeLZ4SingleBlock(in, out []byte, blockSize uint32) (processed uint32, output uint32, err error) {
	var inPos, outPos uint32
	for {
		token := in[inPos]
		inPos++

		length := uint32(token >> 4)
		if length == 15 {
			ln := uint32(in[inPos])
			inPos++
			length += ln
			for ln == 255 {
				ln = uint32(in[inPos])
				inPos++
				length += ln
			}
		}
		if length > 0 {
			m := length / STEPSIZE
			for i := uint32(0); i < m; i++ {
				*(*uint)(unsafe.Pointer(&out[outPos])) = *(*uint)(unsafe.Pointer(&in[inPos]))
				inPos += STEPSIZE
				outPos += STEPSIZE
			}
			m = length % STEPSIZE
			for i := uint32(0); i < m; i++ {
				out[outPos] = in[inPos]
				inPos++
				outPos++
			}
			/*
				copy(out[outPos:outPos+length], in[inPos:inPos+length])
				inPos += length
				outPos += length
			*/
		}

		// last sequence, there is no offset
		if inPos >= blockSize {
			return inPos, outPos, nil
		}

		// two bytes offset, little endian format
		offset := uint32(in[inPos]) | (uint32(in[inPos+1]) << 8)
		inPos += 2
		if offset == 0 || offset > outPos {
			return inPos, outPos, fmt.Errorf("wrong offset %v, outPos: %v, inPos: %v", offset, outPos, inPos)
		}

		matchLength := uint32(token & 0x0F)
		if matchLength == 15 {
			ln := uint32(in[inPos])
			inPos++
			matchLength += ln
			for ln == 255 {
				ln = uint32(in[inPos])
				inPos++
				matchLength += ln
			}
		}
		matchLength += MINMATCH
		overlapCopy(out, outPos, outPos-offset, matchLength)
		outPos += matchLength
	}
}

func EncodeLZ4SingleBlock(in, out []byte) (processed uint32, output uint32, err error) {
	hashTable := make([]uint16, HASHTABLESIZE)
	var ip uint32 = 0
	var op uint32 = 0
	var anchor uint32 = 0
	var val = *(*uint32)(unsafe.Pointer(&in[ip]))
	var h uint32 = (val * GOLDNUMBER) >> (32 - HASHLOG)
	hashTable[h] = uint16(ip)
	mfLimit := uint32(len(in) - MFLIMIT)

	ip++
	val = *(*uint32)(unsafe.Pointer(&in[ip]))
	h = (val * GOLDNUMBER) >> (32 - HASHLOG)

	for {

		var ref uint32
		// find a match
		for {
			if ip > mfLimit {
				goto _last_literals
			}
			ref = uint32(hashTable[h])
			if *(*uint32)(unsafe.Pointer(&in[ref])) == val {
				break
			} else {
				hashTable[h] = uint16(ip)
				ip++
				val = *(*uint32)(unsafe.Pointer(&in[ip]))
				h = (val * GOLDNUMBER) >> (32 - HASHLOG)
			}
		}

		// extend backward
		// TODO can ref be just > 0?
		for ip > anchor && ref > anchor && in[ip-1] == in[ref-1] {
			ip--
			ref--
		}

		token := op
		// encode literal length
		{
			litLength := ip - anchor
			ln := litLength
			if ln >= 15 {
				out[token] = 15 << 4
				op++
				ln -= 15
				for ln >= 255 {
					out[op] = 255
					op++
					ln -= 255
				}
				out[op] = byte(ln) //token
				op++
			} else {
				out[op] = byte(ln) << 4 //token
				op++
			}

			// copy literals
			//fmt.Printf("token from %v, Litlength %v, ip %v, op %v\n", token, litLength, ip, op)
			copy(out[op:op+litLength], in[anchor:anchor+litLength])
			op += litLength
		}

	_next_match:
		// encode offset
		offset := ip - ref
		out[op] = byte(offset)
		op++
		out[op] = byte(offset >> 8)
		op++

		// encode matchLength
		{
			matchLimit := uint32(len(in) - LASTLITERALS)
			matchLength := lz4Count(in, ip+MINMATCH, ref+MINMATCH, matchLimit)
			//fmt.Printf("offset %v op %v, ip %v, match length %v\n", offset, op-2, ip, matchLength)
			//matchLength := dumblz4Count(in, ip+MINMATCH, ref+MINMATCH, matchLimit)
			ip += MINMATCH + matchLength
			if matchLength >= 15 {
				out[token] |= 0x0F
				matchLength -= 15
				for matchLength >= 255 {
					out[op] = 255
					op++
					matchLength -= 255
				}
				out[op] |= byte(matchLength)
				op++
			} else {
				out[token] |= byte(matchLength)
			}
		}
		anchor = ip

		// check end of chunk
		if ip > mfLimit {
			break
		}

		// fill hash table
		// why fill this position?
		val = *(*uint32)(unsafe.Pointer(&in[ip-2]))
		h = (val * GOLDNUMBER) >> (32 - HASHLOG)
		hashTable[h] = uint16(ip - 2)

		// test next position, could be a match, but not consecutive
		// from the one above
		val = *(*uint32)(unsafe.Pointer(&in[ip]))
		h = (val * GOLDNUMBER) >> (32 - HASHLOG)
		ref = uint32(hashTable[h])
		hashTable[h] = uint16(ip)
		if *(*uint32)(unsafe.Pointer(&in[ref])) == val {
			token = op
			out[token] = 0
			op++
			goto _next_match
		}

		// Prepare next loop
		ip++
		val = *(*uint32)(unsafe.Pointer(&in[ip]))
		h = (val * GOLDNUMBER) >> (32 - HASHLOG)
	}
_last_literals:
	// encode last literals
	{
		lastRun := uint32(len(in)) - anchor
		ln := lastRun

		// TODO: check output limit

		if lastRun >= 15 {
			out[op] = 0xF0
			op++
			lastRun -= 15
			for lastRun >= 255 {
				out[op] = 255
				op++
				lastRun -= 255
			}
			out[op] = byte(lastRun)
			op++
		} else {
			out[op] = byte(lastRun << 4)
			op++
		}
		if op+ln > uint32(len(out)) || anchor+ln > uint32(len(in)) {
			fmt.Println("aha", op, anchor, ln, len(out), len(in))
		}
		copy(out[op:op+ln], in[anchor:anchor+ln])

		// end
		return anchor + ln, op + ln, nil
	}
}

func dumblz4Count(in []byte, ip, ref, inLimit uint32) uint32 {
	start := ip
	for ip < inLimit && in[ip+1] == in[ref+1] {
		ip++
		ref++
	}
	return ip - start
}
func lz4Count(in []byte, ip, ref, inLimit uint32) uint32 {
	start := ip
	for ip < inLimit-(STEPSIZE-1) {
		diff := (*(*uint)(unsafe.Pointer(&in[ip]))) ^ (*(*uint)(unsafe.Pointer(&in[ref])))
		if diff == 0 {
			ip += STEPSIZE
			ref += STEPSIZE
			continue
		}
		break
	}
	if _64BITS { //means STEPSIZE is 8
		if ip < (inLimit-3) && *(*uint32)(unsafe.Pointer(&in[ip])) == *(*uint32)(unsafe.Pointer(&in[ref])) {
			ip += 4
			ref += 4
		}
	}
	if (ip < (inLimit - 1)) && *(*uint16)(unsafe.Pointer(&in[ip])) == *(*uint16)(unsafe.Pointer(&in[ref])) {
		ip += 2
		ref += 2
	}
	if ip < inLimit && in[ip] == in[ref] {
		ip++
	}

	return ip - start
}
func CompressBound(in []byte) int {
	return len(in) + ((len(in) / 255) + 16)
}
