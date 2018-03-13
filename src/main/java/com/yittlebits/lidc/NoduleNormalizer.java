package com.yittlebits.lidc;

import java.util.HashMap;
import java.util.Map;
import java.util.Map.Entry;

public class NoduleNormalizer {
    private double distanceThreshold;
    Map<Integer, DoublePoint3> nodules = new HashMap<>();

    public NoduleNormalizer(double distanceThreshold) {
	this.distanceThreshold = distanceThreshold;
    }

    public int getNormalizedId(DoublePoint3 pt) {
	// Find the closest
	double minDistance = 10e10;
	int matchIdx = -1;
	for (Entry<Integer, DoublePoint3> nodule : nodules.entrySet()) {
	    double distance = nodule.getValue().distanceTo(pt);
	    if (distance < minDistance) {
		matchIdx = nodule.getKey();
		minDistance = distance;
	    }
	}
	if (minDistance > distanceThreshold || matchIdx == -1) {
	    // Insert first one
	    matchIdx = nodules.size() + 1;
	    nodules.put(matchIdx, pt);

	}
	return matchIdx;
    }
}
