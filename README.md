**Format Image**

Service to reduce a image

Input:

    file - image
    string - maxWidth
    string - maxHeight
    file - watermark @todo

Process:

    Downscales an image preserving its aspect ratio to the maximum dimensions (maxWidth, maxHeight). It will return the original image if original sizes are smaller than the provided dimensions.

Output:

    file - reduced image

Run the http server:

    go run format-image.go

Example with cURL:

    curl --location 'http://localhost:3000/reduce' \
    --form 'image=@"/path/to/image.jpg"' \
    --form 'maxWidth="500"' \
    --form 'maxHeight="250"'