import java.io.*;
import java.net.*;

class ChatClient {

  public static void main(String... args) throws IOException {
    var host = "localhost";
    var port = 6666;

    try (
        var s = new Socket(host, port);
        var stdin = new BufferedReader(new InputStreamReader(System.in));
        var out = new PrintWriter(new OutputStreamWriter(s.getOutputStream()), true);
        
    ) {
      new ServerReaderThread(s).start();
      String line;

      while (true) {
        System.out.print("> ");
        if ((line = stdin.readLine()) == null) {
          break;
        }
        out.println(line);
      }
    }
  }
}

class ServerReaderThread extends Thread {
    private Socket s;

  ServerReaderThread(Socket s){
      this.s = s;

  }
  public void run() {
    try (
      var in = new BufferedReader(new InputStreamReader(s.getInputStream()));
    ) {
      String response;
      while (true) {
        if ((response = in.readLine()) == null) {
          break;
        }
        System.out.println("SERVER > " + response);
        System.out.print("> ");

      }
    } catch(Exception e) {
      e.printStackTrace();
    }
  }
}
//reading from socket and keyboard is 2 different threads.
// there must be 2 threads
