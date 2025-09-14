import easyocr
import os

def perform_ocr(image_path):
    if not os.path.exists(image_path):
        return f"Error: File not found at the specified path: {image_path}"

    try:
        reader = easyocr.Reader(['pl', 'en'], gpu=False)
        result = reader.readtext(image_path, paragraph=True)
        extracted_text = "\n".join([item[1] for item in result])

        return extracted_text

    except Exception as e:
        return f"An unexpected error occurred: {e}"

if __name__ == "__main__":
    path_to_image = "img.png"
    text_from_image = perform_ocr(path_to_image)
    print("Recognized text: " + text_from_image)