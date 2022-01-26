package com.example.client;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestTemplate;

import java.nio.charset.Charset;
import java.nio.charset.StandardCharsets;
import java.security.SecureRandom;
import java.util.Random;

public class Followers implements Runnable {

    private Random random =  new Random();
    private RestTemplate restTemplate = new RestTemplate();
    HttpHeaders headers = new HttpHeaders();

    final String URI_TEST = "http://localhost:8080/api/user/followTest";
    String easy = RandomString.digits + "ACEFGHJKLMNPQRUVWXYabcdefhijkprstuvwx";

    public void init(){
        //restTemplate.getMessageConverters().add(new MappingJackson2HttpMessageConverter());
        restTemplate.getMessageConverters().add(0, new MappingJackson2HttpMessageConverter());
//        headers.setContentType(MediaType.APPLICATION_JSON);
//        headers.clearContentHeaders();
        headers.clear();
        headers.add("Content-Type", "application/json;charset=UTF-8");
//        headers.add(HttpHeaders.ACCEPT, MediaType.APPLICATION_JSON_VALUE);
//        headers.add(HttpHeaders.ACCEPT_CHARSET, StandardCharsets.UTF_8.name());
        System.out.println(headers.toString());
    }

    @Override
    public void run() {
        this.init();
        for (int k=0; k< 1000; k++){
            FollowerDto followerDto = new FollowerDto();
            int i = random.nextInt(999)+1;
            int j = random.nextInt(999)+1;
            followerDto.setMyId(String.format("%d", i));
            followerDto.setFriendId(String.format("%d", j));
            try {
                if (i != j) {
                    HttpEntity<FollowerDto> request = new HttpEntity<>(followerDto, headers);
                    StatusResponse addedPost = restTemplate.postForObject(URI_TEST, followerDto, StatusResponse.class);
                    System.out.println("follower added : " + addedPost.getStatus() + "   " + followerDto.getMyId() + " " + followerDto.getFriendId());
                }
            }
            catch (Exception e){
                continue;
            }
        }

    }
}