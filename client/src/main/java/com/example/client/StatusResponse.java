package com.example.client;

import com.fasterxml.jackson.annotation.JsonInclude;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
public class StatusResponse {
    @JsonProperty(value ="status")
    @JsonInclude(JsonInclude.Include.NON_NULL)
    private String status;
}


