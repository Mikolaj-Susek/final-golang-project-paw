import grpc
from concurrent import futures
import ocr_pb2
import ocr_pb2_grpc
import easyocr
import logging

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')

class OcrServiceImpl(ocr_pb2_grpc.OcrServiceServicer):
    def __init__(self):
        try:
            self.reader = easyocr.Reader(['pl', 'en'], gpu=False)
        except Exception as e:
            logging.error(f"Failed to initialize EasyOCR: {e}")
            raise

    def PerformOcr(self, request, context):
        try:
            logging.info("Received new OCR request.")
            image_bytes = request.image_data

            result = self.reader.readtext(image_bytes, paragraph=True)

            extracted_text = "\n".join([item[1] for item in result])

            logging.info(f"OCR operation completed successfully. Text length: {len(extracted_text)} characters.")

            return ocr_pb2.OcrResponse(extracted_text=extracted_text)

        except Exception as e:
            error_message = f"An unexpected error occurred during OCR processing: {e}"
            logging.error(error_message)
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details(error_message)
            return ocr_pb2.OcrResponse()

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    ocr_pb2_grpc.add_OcrServiceServicer_to_server(OcrServiceImpl(), server)

    server.add_insecure_port('[::]:50051')

    logging.info("Starting server on port 50051...")
    server.start()
    logging.info("Server has started and is waiting for connections.")

    server.wait_for_termination()

if __name__ == '__main__':
    serve()
    