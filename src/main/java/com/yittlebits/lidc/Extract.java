package com.yittlebits.lidc;

import java.util.logging.Logger;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import javax.xml.xpath.XPath;
import javax.xml.xpath.XPathConstants;
import javax.xml.xpath.XPathFactory;

import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.DefaultParser;
import org.apache.commons.cli.HelpFormatter;
import org.apache.commons.cli.Options;
import org.w3c.dom.Document;
import org.w3c.dom.Node;

class Extract {
  static Logger logger = Logger.getLogger(Extract.class.getName());

  static Options options;
  //@formatter:off
  static String usage = String.join("\n",
      "Extract <command> <xml>",
      "",
      "Extract information from LIDC XMl and DICOM files",
      "Commands are:",
      "  SeriesInstanceUID",
      "    Extract the SeriesInstanceUID from the XML file (/LidcReadMessage/ResponseHeader/SeriesInstanceUid tog)",
      "    print it and exit.",
      "",
      "  segment <xml> <DICOM Input> <DICOM output>",
      "    Create new DICOM objects indicating the position of nodules from the XML",
      "    Parameters are DICOM input and output directories"
   );
  //@formatter:on

  static void printNode(Node node) {
    logger.info(node.getNodeName() + ": " + node.getNodeValue());
    logger.info("Children: " + node.getChildNodes());

  }

  public static void printUsageAndDie() {
    HelpFormatter formatter = new HelpFormatter();
    formatter.printHelp(usage, options);
    System.exit(1);
  }

  public static void main(String[] args) throws Exception {

    // create Options object
    options = new Options();
    CommandLine cl = new DefaultParser().parse(options, args);

    if (cl.getArgList().size() < 1) {
      printUsageAndDie();
    }

    String command = cl.getArgList().get(0);
    String xmlFile = cl.getArgList().get(1);

    DocumentBuilderFactory factory = DocumentBuilderFactory.newInstance();
    factory.setNamespaceAware(false);
    DocumentBuilder builder;
    Document doc = null;
    builder = factory.newDocumentBuilder();
    doc = builder.parse(xmlFile);

    // create an XPathFactory
    XPathFactory xFactory = XPathFactory.newInstance();

    // create an XPath object
    XPath xpath = xFactory.newXPath();

    if (command.toLowerCase().equals("seriesinstanceuid")) {
      // get the SeriesInstanceUid
      String seriesInstanceUID = (String) xpath.compile("/LidcReadMessage/ResponseHeader/SeriesInstanceUid/text()").evaluate(doc, XPathConstants.STRING);
      System.out.println(seriesInstanceUID);
      System.exit(0);
    }

    if (command.toLowerCase().equals("segment")) {
      Segmenter segmenter = new Segmenter();
      segmenter.segment(cl);
    }

  }
}