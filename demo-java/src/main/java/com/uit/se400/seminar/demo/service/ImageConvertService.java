package com.uit.se400.seminar.demo.service;

import org.springframework.stereotype.Service;
import org.springframework.web.multipart.MultipartFile;

import javax.imageio.ImageIO;
import java.awt.image.BufferedImage;
import java.io.*;

@Service
public class ImageConvertService {

    public static byte[] convertToMonochrome(MultipartFile file) {
        try {
            // Read image to buffer
            InputStream inputStream = new ByteArrayInputStream(file.getBytes());
            BufferedImage image = ImageIO.read(inputStream);
            int width = image.getWidth();
            int height = image.getHeight();

            // Create result image object
            BufferedImage monochromeImage = new BufferedImage(width, height, BufferedImage.TYPE_BYTE_GRAY);

            // Convert each pixel to monochrome
            for (int y = 0; y < height; y++) {
                for (int x = 0; x < width; x++) {
                    int rgb = image.getRGB(x, y);
                    int r = (rgb >> 16) & 0xFF;
                    int g = (rgb >> 8) & 0xFF;
                    int b = rgb & 0xFF;

                    // Calculate grayscale
                    // Luminance formula: Gray = 0.299×R + 0.587×G + 0.114×B
                    int gray = (int) (0.299 * r + 0.587 * g + 0.114 * b);

                    // Set grayscale value into pixels
                    int newRgb = (gray << 16) | (gray << 8) | gray;
                    monochromeImage.setRGB(x, y, newRgb);
                }
            }

            // Write output image
            ImageIO.write(monochromeImage, "jpg", new File("output.jpg"));
            System.out.println("Conversion completed.");
            ByteArrayOutputStream byteArrayOutputStream = new ByteArrayOutputStream();
            ImageIO.write(monochromeImage, "jpg", byteArrayOutputStream);
            return byteArrayOutputStream.toByteArray();
        } catch (IOException e) {
            throw new RuntimeException(e);
        }
    }
}
