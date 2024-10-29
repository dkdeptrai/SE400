package com.uit.se400.seminar.demo.repository;

import com.uit.se400.seminar.demo.entity.Order;
import org.springframework.data.jpa.repository.JpaRepository;

public interface OrderRepository extends JpaRepository<Order, Long> {
}
