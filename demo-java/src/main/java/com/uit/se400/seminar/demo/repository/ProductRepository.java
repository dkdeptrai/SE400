package com.uit.se400.seminar.demo.repository;

import com.uit.se400.seminar.demo.entity.Product;
import org.springframework.data.jpa.repository.JpaRepository;

public interface ProductRepository extends JpaRepository<Product, Long> {
}