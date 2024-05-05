package vector

import (
	"errors"
	"fmt"
    "os"
    "math"
    "bytes"
    "bufio"
    "encoding/gob"
)

func (db *DB) Insert(v Vector) error {
	if len(v.Data) < db.Dim {
		return errors.New("Input vector dimension doesn't match the database dimenions")

	}
	if len(db.clusters) == 0 {
        fmt.Println("No clusters present")
        if err := db.addCluster(v); err != nil {
            return err
        }
		return nil
	}
    
    var match *cluster = nil
    matchInd := 0
    var minD float64 = math.Inf(1)
    for i, c := range db.clusters {
        if val, err := db.Sim.calc(c.center ,v); err != nil {
            return err
        } else {
            fmt.Println("Distance Value: ", val, " ; Threshold: ", db.Sim.threshold());
            if val < db.Sim.threshold() {
                if (match == nil) || (val < minD) {
                    match = &c
                    minD = val
                    matchInd = i
                }
            }
        }
    }
    if match == nil {
        if err := db.addCluster(v); err != nil {
            return err
        }
    } else {
        if err := match.fileAppend(v); err != nil {
            return err 
        }
       match.adjustCenter(v)
       db.clusters[matchInd] = *match
    }
	return nil
}


func (c *cluster) adjustCenter(v Vector){
    vec := c.center
    for i := range len(vec.Data) {
        vec.Data[i] *= (float64)(c.pts - 1)
        vec.Data[i] += v.Data[i]
        vec.Data[i] /= (float64)(c.pts)
    }
    fmt.Printf("New Center of Cluster %d is %v\n", c.id, vec);
    c.center = vec
}

func newFile(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, errors.New("Failed to create a new cluster file")
	}
	defer f.Close()
	return f, nil
}

func openInAppend(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
}

func (c *cluster) fileRead() ([]Vector, error){
    //add caching to keep recent reads in memory
    file, err := os.Open(c.fptr.path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    // Create a new scanner for the file.
    scanner := bufio.NewScanner(file)

    // Create a slice to store the vectors.
    vectors := []Vector{}

    // Read the file line by line.
    for scanner.Scan() {
        // Decode the JSON line into a Person struct.
        var vector Vector
        decoder := gob.NewDecoder(bytes.NewReader(scanner.Bytes()))
        err = decoder.Decode(&vector)

        if err != nil {
            return nil, err
        }
        vectors = append(vectors,vector)
    }
     return vectors, nil
}

func (c *cluster)fileAppend(v Vector) ( error) {
    buf := new(bytes.Buffer)

    // Create a new gob encoder.
    encoder := gob.NewEncoder(buf)

    // Encode the vector struct to the byte buffer.
    err := encoder.Encode(v)
    if err != nil {
        return err
    }
    buf.Write([]byte{'\n'})
    if _, err := c.fptr.ptr.Write(buf.Bytes()); err != nil {
        return err  
    } else {
        c.pts += 1
        //fmt.Printf("Wrote %d bytes to file\n", n)
        return nil
    }
}

func (db *DB) addCluster(v Vector) (error) {
    if cl, err := newCluster(len(db.clusters)+1); err != nil {
        return err
    } else {
        cl.center = v
        cl.fileAppend(v)
        db.clusters = append(db.clusters, cl)
    }
    return nil
}

func newCluster(num int) (cluster, error) {
	path := fmt.Sprintf("%s/%d.dat", ddir, num)
	if _, err := newFile(path); err != nil {
		return cluster{}, err
	}

	f, err := openInAppend(path)
	if err != nil {
		return cluster{}, err
	}

    return cluster{id: num, fptr: file{path: path, ptr: f}, pts : 0}, nil

}

func (db *DB) Search(v Vector) (Vector, error) {
    var match *cluster = nil
    var minD float64 = math.Inf(1)
    for _, c := range db.clusters {
        if val, err := db.Sim.calc(c.center ,v); err != nil {
            return Vector{}, err
        } else {
            fmt.Println("Distance Value: ", val, " ; Threshold: ", db.Sim.threshold());
            if val < db.Sim.threshold() {
                if (match == nil) || val < minD {
                    match = &c
                    minD = val
                }
            }
        }
    }
    if match == nil {
        return Vector{}, errors.New("No match found")
    }
        
    if vecs, err := match.fileRead(); err != nil {
        return Vector{}, err
    } else {
        counts := map[int]int{}
        for _, vec:= range vecs {
            counts[vec.Type] += 1
        }
        maxD := -1
        maxT := 0
        for k,v := range counts {
            if maxD < v {
                maxD = v
                maxT = k
            }
        }
        fmt.Printf("The type of cluster is %d\n", maxT);
    }
 
    return match.center, nil
}

// de-bulk index: decide if an index is bulky and move items around
