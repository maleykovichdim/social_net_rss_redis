package com.example.client;

import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.converter.json.MappingJackson2HttpMessageConverter;
import org.springframework.web.client.RestTemplate;
import org.springframework.web.util.UriComponents;
import org.springframework.web.util.UriComponentsBuilder;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Random;

public class Rss implements Runnable {

    private final int NUM_IN_REDIS = 1000;

    private Random random =  new Random();
    private RestTemplate restTemplate = new RestTemplate();
    HttpHeaders headers = new HttpHeaders();

    final String URI_TEST = "http://localhost:8080/api/user/rss_feed";
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
            //FollowerDto followerDto = new FollowerDto();
            int i = random.nextInt(999)+1;
//            int i = 441;

            UriComponentsBuilder builder = UriComponentsBuilder.fromHttpUrl(URI_TEST)
                    .queryParam("my_id",String.format("%d", i));



            System.out.println("posts ----------------------------------------------------"+builder.toUriString());
            PostDto[] posts = restTemplate.getForObject(builder.toUriString(), PostDto[].class);
            if (posts != null) {
                for (PostDto post : posts) {
                    System.out.println(post.toString());
                }
            }
            System.out.println("----------------------------------------------------------");
//                System.out.println("follower added : " + addedPost.getStatus() + "   " + followerDto.getMyId() + " " + followerDto.getFriendId());

        }

    }
}