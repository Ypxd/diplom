import UIKit

struct HTTPResponse: Decodable {
    let error_text: String
    let has_error: Bool
    let message: String
}

class LogInViewController: UIViewController {
    
    @IBOutlet weak var errorLable: UILabel!
    @IBOutlet weak var password: UITextField!
    @IBOutlet weak var login: UITextField!
    @IBOutlet weak var logInButton: UIButton!
    @IBOutlet weak var singIn: UIButton!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    
    @IBAction func singInButton(_ sender: UIButton) {
    }
    
    @IBAction func logInButton(_ sender: UIButton) {
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/") else { return }
        var request = URLRequest(url: url)
        let param = ["login": login.text!, "password": password.text!]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        
                        DispatchQueue.main.async {
                            self.dismiss(animated: true, completion: nil)
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        self.errorLable.text = httpResp.error_text
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
    }
}
