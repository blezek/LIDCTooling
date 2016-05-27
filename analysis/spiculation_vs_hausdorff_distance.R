png("spiculation_vs_hausdorff_distance.png", height = 1200, width = 1600)
d <- read.csv("spiculation_vs_hausdorff_distance.csv")
plot ( d$spiculation, d$average_hausdorff_distance, main="Segmentation distance vs. Spiculation", xlab="Spiculation (1-5)", ylab="Avg. Hausdorff Distance (mm)")
dev.off()

