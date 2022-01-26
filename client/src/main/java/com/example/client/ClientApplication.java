package com.example.client;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class ClientApplication {

    public static void main(String[] args) {
        SpringApplication.run(ClientApplication.class, args);
        Thread threadPost = new Thread(new Posts());
        Thread threadFollower = new Thread(new Followers());
        Thread threadRss = new Thread(new Rss());
       threadPost.start();
        threadFollower.start();
        threadRss.start();
    }
}
