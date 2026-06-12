import grpc
from concurrent import futures

import accountant_pb2
import accountant_pb2_grpc
from mailSender import SendMessageSMTP


class AccountantService(
    accountant_pb2_grpc.AccountantServiceServicer
):
    def SendAccountant(self, request, context):
        print(request.message)

        return accountant_pb2.SendAccountantResponse(
            success=True,
            error=""
        )


def serve():
    server = grpc.server(
        futures.ThreadPoolExecutor(max_workers=10)
    )

    accountant_pb2_grpc.add_AccountantServiceServicer_to_server(
        AccountantService(),
        server,
    )

    server.add_insecure_port("[::]:50051")
    server.start()
    server.wait_for_termination()


if __name__ == "__main__":
    SendMessageSMTP("Тестовое")
    serve()