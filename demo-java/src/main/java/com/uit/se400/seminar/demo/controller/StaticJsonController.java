package com.uit.se400.seminar.demo.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/static-json")
public class StaticJsonController {
    @GetMapping
    public Map<String, Object> getStaticJson() {
        Map<String, Object> response = new HashMap<>();
        response.put("message", "Hello, World!");
        response.put("status", "success");

        Map<String, Object> user = new HashMap<>();
        user.put("id", 1);
        user.put("name", "John Doe");
        user.put("email", "john.doe@example.com");

        Map<String, Object> post1 = new HashMap<>();
        post1.put("id", 101);
        post1.put("title", "First Post");
        post1.put("content", "This is the content of the first post.");

        Map<String, Object> post2 = new HashMap<>();
        post2.put("id", 102);
        post2.put("title", "Second Post");
        post2.put("content", "This is the content of the second post.");

        response.put("data", Map.of(
                "user", user,
                "posts", List.of(post1, post2)
        ));

        return response;
    }
}
