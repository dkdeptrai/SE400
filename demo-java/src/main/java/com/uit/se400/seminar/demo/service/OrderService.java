package com.uit.se400.seminar.demo.service;

import com.uit.se400.seminar.demo.entity.Order;
import com.uit.se400.seminar.demo.entity.Product;
import com.uit.se400.seminar.demo.repository.OrderRepository;
import com.uit.se400.seminar.demo.repository.ProductRepository;
import jakarta.transaction.Transactional;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class OrderService {

    private final ProductRepository productRepository;

    private final OrderRepository orderRepository;

    @Autowired
    public OrderService(ProductRepository productRepository, OrderRepository orderRepository) {
        this.productRepository = productRepository;
        this.orderRepository = orderRepository;
    }

    @Transactional
    public boolean checkStockAndCreateOrder(Order order) throws Exception {
        Product product = productRepository.findById(order.getProduct().getId())
                .orElseThrow(() -> new Exception("Product does not exist"));

        if (product.getStock() < order.getQuantity()) {
            return false;
        }

        // Update stock
        product.setStock(product.getStock() - order.getQuantity());
        productRepository.save(product);

        // Calculate total price and set status
        order.setTotalPrice(product.getPrice() * order.getQuantity());
        order.setStatus("Pending");
        orderRepository.save(order);

        return true;
    }

    public Order getOrder(Long orderId) throws Exception {
        return orderRepository.findById(orderId)
                .orElseThrow(() -> new Exception("Order not found"));
    }
}
