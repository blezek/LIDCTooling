package com.yittlebits.lidc.fill;

import com.yittlebits.lidc.Point3;

public class Edge {
  public Edge next = null;

  public int yUpper;
  public int dx;
  public int dy;
  public int dy2;
  public int dx2;
  public int dydx2;
  public int r;
  public int xInc;
  public int x;

  public static void insertEdge(Edge list, Edge edge) {
    Edge p;
    Edge q = list;

    p = q.next;
    while (p != null) {
      if (edge.x < p.x) {
        p = null;
      } else {
        q = p;
        p = p.next;
      }
    }
    edge.next = q.next;
    q.next = edge;

  }

  // ------------------------------------------------------------------------------
  // MarkEdgeRec
  //
  // This routine stores lower-y coordinate and inverse slope for each slope.
  //
  // Adjust and store upper-y coordinate for edges that are the
  // lower member of a monotonically increasing or decreasing
  // pair of edges.
  // ------------------------------------------------------------------------------
  public static void makeEdgeRec(int x1, int y1, int x2, int y2, Edge edge, Edge edges[]) {
    // p1 is lower than p2
    edge.dx = Math.abs(x2 - x1);
    edge.dy = Math.abs(y2 - y1);
    edge.dx2 = edge.dx << 1;
    edge.dy2 = edge.dy << 1;
    if (x1 < x2) {
      edge.xInc = 1;
    } else {
      edge.xInc = -1;
    }
    edge.x = x1;

    // < 45 degree slope
    if (edge.dy <= edge.dx) {
      edge.dydx2 = (edge.dy - edge.dx) << 1;
      edge.r = edge.dy2 - edge.dx;
    }
    // > 45 degree slope
    else {
      edge.dydx2 = (edge.dx - edge.dy) << 1;
      edge.r = edge.dx2 - edge.dy;
    }

    edge.yUpper = y2;

    Edge.insertEdge(edges[y1], edge);
  }

  // ------------------------------------------------------------------------------
  // BuildEdgeList
  // ------------------------------------------------------------------------------
  public static void buildEdgeList(Point3 pPnts[], Edge edges[]) {
    int i, x1, x2, y1, y2, nPnts = pPnts.length;
    Edge edge = null;

    // Previous point (which is the last point for the first point):
    i = nPnts - 1;
    x1 = pPnts[i].x;
    y1 = pPnts[i].y;

    // Foreach point:
    for (i = 0; i < nPnts; i++) {
      // Convert coordinates from Extent-based IJK to 0-based:
      x2 = pPnts[i].x;
      y2 = pPnts[i].y;

      if (y1 != y2) {
        // This is a non-horizontal line.
        edge = new Edge();
        if (y1 < y2) {
          // This is an up-going edge.
          makeEdgeRec(x1, y1, x2, y2, edge, edges);
        } else {
          // This is a down-going edge.
          makeEdgeRec(x2, y2, x1, y1, edge, edges);
        }
      }
      // Advance:
      x1 = x2;
      y1 = y2;
    }
  }

}
