package zstrset

import "strings"

type zstrset struct {
	// concat []string tags into a compact format. Equivalent to:
	//     strings.Join(tags, "")
	data []byte

	// meta stores size and the cached location of `:`
	//     meta[i*2]   = size of tag[i]
	//     meta[i*2+1] = strings.IndexByte(tag[i], ':')
	meta []int
}

func (z *zstrset) each(f func(group, val []byte)) {
	offset := 0
	for i := 0; i < len(z.meta); i += 2 {
		size := z.meta[i]
		tag := z.data[offset : offset+size]

		if sepidx := z.meta[i+1]; sepidx < 0 {
			f(nil, tag)
		} else {
			f(tag[:sepidx], tag[sepidx+1:])
		}
		offset += size
	}
}

func makezstrset(tags []string) zstrset {
	nbyte := 0

	meta := make([]int, len(tags)*2)

	for i, t := range tags {
		nbyte += len(t)

		meta[i*2] = len(t)
		meta[i*2+1] = strings.IndexByte(t, ':')
	}

	// 2nd loop to avoid expensive realloc
	data := make([]byte, 0, nbyte)
	for _, t := range tags {
		data = append(data, t...)
	}

	return zstrset{
		data: data,
		meta: meta,
	}
}
