import UIKit

class RegisterViewController: UIViewController {

    @IBOutlet weak var errorLabel: UILabel!
    @IBOutlet weak var singInButton: UIButton!
    @IBOutlet weak var logInButton: UIButton!
    @IBOutlet weak var age: UITextField!
    @IBOutlet weak var name: UITextField!
    @IBOutlet weak var email: UITextField!
    @IBOutlet weak var password: UITextField!
    @IBOutlet weak var login: UITextField!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
    
    @IBAction func logInButton(_ sender: UIButton) {
    }
    
    @IBAction func singInButton(_ sender: UIButton) {
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/register") else { return }
        var request = URLRequest(url: url)
        let param = ["login": login.text!, "password": password.text!, "email": email.text!, "name": name.text!, "age": age.text!]
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
                        DispatchQueue.main.async {
                            self.dismiss(animated: true, completion: nil)
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        self.errorLabel.text = httpResp.error_text
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    
}
