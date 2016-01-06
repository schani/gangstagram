package main

import (
    "bufio"
    "fmt"
    "math"
    "os"
    "strconv"
)

type cluster struct {
    mean float64
    sumOfSquares float64
    count int64
}

type clusterer struct {
    maxK int
    clusters []cluster
}

func (c *cluster) distance(x float64) float64 {
    return math.Abs(x - c.mean)
}

func (c *cluster) update(x float64) {
    oldMean := c.mean
    c.mean = oldMean + (x - oldMean) / float64(c.count + 1)
    c.sumOfSquares += x * x
    c.count++
}

func positiveCount (i int64) float64 {
    //return float64(i)
    return math.Max(float64(i), 1)
}

func (c *cluster) combine(c2 *cluster) cluster {
    newCount := c.count + c2.count
    return cluster{
        mean: (c.mean * float64(c.count) + c2.mean * float64(c2.count)) / positiveCount(newCount),
        sumOfSquares: c.sumOfSquares + c2.sumOfSquares,
        count: newCount,
    }
}

func (c *cluster) variance() float64 {
    if c.count < 1 {
        return 0
    }
    return c.sumOfSquares / float64(c.count) - c.mean * c.mean
}

// a < b
func (cr *clusterer) closestClusters() (a, b int) {
    k := len(cr.clusters)
    if k < 2 {
        panic("can't find closest clusters with less than 2")
    }

    minDist := float64(-1)
    
    for i := 0; i < k; i++ {
        for j := i + 1; j < k; j++ {
            dist := cr.clusters[i].distance(cr.clusters[j].mean)
            if minDist < 0 || dist < minDist {
                a = i
                b = j
                minDist = dist
            }
        }
    }
    
    if minDist < 0 || a >= b {
        panic("we're stupid")
    }
    return a, b
}

func (cr *clusterer) add(x float64) {
    k := len(cr.clusters)
    if k != 0 {
        var i int
        minI := 0
        minDistance := cr.clusters[minI].distance(x) 
        for i = 1; i < len(cr.clusters); i++ {
            dist := cr.clusters[i].distance(x)
            if dist < minDistance {
                minI = i
                minDistance = dist
            }
        }
        
        cr.clusters[minI].update(x)
        
        if k == cr.maxK {
            a, b := cr.closestClusters()
            c := cr.clusters[a].combine(&cr.clusters[b])
            // remove a, b
            cr.clusters = append(cr.clusters[0:a], append(cr.clusters[a+1:b], cr.clusters[b+1:k]...)...)
            cr.clusters = append(cr.clusters, c)
        } 
    }
    
    cr.clusters = append(cr.clusters, cluster{mean: x, sumOfSquares: 0, count: 0})
}

func main() {
    cr := clusterer{maxK: 10, clusters: []cluster{}}
    //c := cluster{mean: 10, sumOfSquares: 0, count: 0}
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        line := scanner.Text()
        x, err := strconv.ParseFloat(line, 64)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Error: Couldn't parse line `%s': %s", line, err.Error())
            continue
        }
        //c = c.combine(&cluster{mean: x, sumOfSquares: 0, count: 1})
        cr.add(x)
        //c.update(x)
    }
    
    for _, c := range(cr.clusters) {
        fmt.Printf("mean: %f  stddev: %f  count: %d\n", c.mean, math.Sqrt(c.variance()), c.count)
    }
}
