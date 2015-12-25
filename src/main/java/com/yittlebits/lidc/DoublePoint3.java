package com.yittlebits.lidc;

public class DoublePoint3 {

  public double x, y, z;

  public DoublePoint3(double xx, double yy, double zz) {
    x = xx;
    y = yy;
    z = zz;
  }

  public double distanceTo(DoublePoint3 pt) {
    double dx = x - pt.x;
    double dy = y - pt.y;
    double dz = z - pt.z;
    return Math.sqrt(dx * dx + dy * dy + dz * dz);
  }

}
