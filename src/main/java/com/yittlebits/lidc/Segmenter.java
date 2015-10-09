package com.yittlebits.lidc;

import java.io.File;
import java.io.IOException;
import java.nio.file.FileVisitResult;
import java.nio.file.FileVisitor;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.List;
import java.util.Map;
import java.util.logging.Logger;
import java.util.stream.Collectors;

import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.xpath.XPath;
import javax.xml.xpath.XPathFactory;

import org.apache.commons.cli.CommandLine;
import org.dcm4che2.data.DicomObject;
import org.dcm4che2.data.Tag;
import org.w3c.dom.Document;

public class Segmenter {
  static Logger logger = Logger.getLogger(Segmenter.class.getName());

  Document doc = null;
  XPath xpath;
  File inputDirectory;
  File outputDirectory;

  // Sorted by Tag.SliceLocation
  List<DicomObject> dicomObjects;

  // Indexed by SOPInstanceUID
  Map<String, DicomObject> uidToInstance;

  public void segment(CommandLine cl) throws Exception {
    if (cl.getArgList().size() < 4) {
      Extract.printUsageAndDie();
    }

    String xmlFile = cl.getArgList().get(1);
    inputDirectory = new File(cl.getArgs()[2]);
    if (!inputDirectory.exists() || !inputDirectory.isDirectory()) {
      System.err.println("input directory " + inputDirectory + " must exist and be a directory");
      Extract.printUsageAndDie();
    }

    outputDirectory = new File(cl.getArgs()[3]);
    if (!outputDirectory.exists() || !outputDirectory.isDirectory()) {
      System.err.println("output directory " + outputDirectory + " must exist and be a directory");
      Extract.printUsageAndDie();
    }

    DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
    factory.setNamespaceAware(false);
    doc = factory.newDocumentBuilder().parse(xmlFile);
    XPathFactory xFactory = XPathFactory.newInstance();
    xpath = xFactory.newXPath();

    // Load DICOM
    loadDICOM();
    logger.info("Loaded " + dicomObjects.size() + " DICOM from " + inputDirectory);
    for (DicomObject dicom : dicomObjects) {
      logger.info("Slice location: " + dicom.getDouble(Tag.SliceLocation));
    }

  }

  private void loadDICOM() throws Exception {
    dicomObjects = Files.walk(inputDirectory.toPath()).filter(p -> {
      return !p.toFile().isDirectory();
    }).map(p -> {
      try {
        return TagLoader.loadTags(p.toFile(), true);
      } catch (Exception e) {
        return null;
      }
    }).sorted((d1, d2) -> {
      return Double.compare(d1.getDouble(Tag.SliceLocation, -1000000.0), d2.getDouble(Tag.SliceLocation, -100000.0));
    }).collect(Collectors.toList());

    // index by uid, how easy is that!
    uidToInstance = dicomObjects.stream().collect(Collectors.toMap(item -> {
      return item.getString(Tag.SOPInstanceUID);
    }, item -> item));

  }
}
