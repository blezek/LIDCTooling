package com.yittlebits.lidc.fill;

import com.yittlebits.lidc.Point3;

import niftijio.NiftiVolume;

public class Fill {

  // volume.data.set((int) x, (int) y, z, 0, labelValue);

  public static void drawPolygon(NiftiVolume volume, Point3 pPnts[], int labelValue, int z) {
    int nx = volume.data.sizeX();
    int ny = volume.data.sizeY();

    int i, scan, done, x, x1, x2;
    Edge edges[] = new Edge[2048], active = null, p = null, q = null, del = null;

    // Build a list of edges for each pixel in the Y direction:
    for (i = 0; i < ny; i++) {
      edges[i] = new Edge();
    }
    Edge.buildEdgeList(pPnts, edges);
    active = new Edge();

    // Careful: scan is a 0-based Y-coordinate, not Extent-based
    for (scan = 0; scan < ny; scan++) {
      // BuildActiveList(int scan, Edge *active, Edge *edges[])
      p = edges[scan].next;
      while (p != null) {
        q = p.next;
        Edge.insertEdge(active, p);
        p = q;
      }

      if (active.next != null) {
        // DeleteFromActiveList(int scan, Edge *active)
        q = active;
        p = active.next;
        while (p != null) {
          if (scan >= p.yUpper) {
            p = p.next;

            // Delete After q
            del = q.next;
            q.next = del.next;
          } else {
            q = p;
            p = p.next;
          }
        }

        // Fill Scan:
        p = active.next;
        while (p != null) {
          q = p.next;
          x1 = p.x;
          x2 = q.x;

          // Fill from left edge up to, but not including, the right edge:
          for (x = x1; x < x2; x++) {
            // nx, y
            volume.data.set(x, scan, z, 0, labelValue);
            // ptr[x] = 1;
          }
          p = q.next;
        }

        // Update Active List:
        q = active;
        p = active.next;
        while (p != null) {
          // < 45 degree slope
          if (p.dy <= p.dx) {
            done = 0;
            while (done == 0) {
              p.x += p.xInc;
              if (p.r <= 0)
                p.r += p.dy2;
              else {
                done = 1;
                p.r += p.dydx2;
              }
            }
          }
          // > 45
          else {
            if (p.r <= 0) {
              p.r += p.dx2;
            } else {
              p.x += p.xInc;
              p.r += p.dydx2;
            }
          }
          q = p;
          p = p.next;
        }

        // Resort Active List:
        p = active.next;
        active.next = null;
        while (p != null) {
          q = p.next;
          Edge.insertEdge(active, p);
          p = q;
        }
      } // if(active.next)
    } // for
  }

}
